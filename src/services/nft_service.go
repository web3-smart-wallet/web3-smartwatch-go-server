package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/web3-smart-wallet/src/api"
)

type NFTService struct {
	apiURL string
}

type NFTServiceInterface interface {
	GetNFTs(address string, includeMetadata bool, pageToken string) ([]api.NFT, string, error)
}

func NewNFTService() NFTServiceInterface {
	apiURL := os.Getenv("ANKR_API_URL")
	if apiURL == "" {
		panic("ANKR_API_URL environment variable is not set")
	}

	return &NFTService{
		apiURL: apiURL,
	}
}

func (s *NFTService) GetNFTs(address string, includeMetadata bool, pageToken string) ([]api.NFT, string, error) {
	// 构建请求体
	params := map[string]interface{}{
		"blockchain":      "base",
		"walletAddress":   address,
		"includeMetadata": includeMetadata,
	}

	// 如果有 pageToken，添加到请求参数中
	if pageToken != "" {
		params["pageToken"] = pageToken
	}

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "ankr_getNFTsByOwner",
		"params":  params,
		"id":      1,
	}

	// 发送请求
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request: %v", err)
	}

	fmt.Printf("Sending request to %s with payload: %s\n", s.apiURL, string(jsonData))

	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时时间
	}
	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch NFTs: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response from Ankr API: %s\n", string(respBody))

	// 解析响应
	var response struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  struct {
			NextPageToken string `json:"nextPageToken"`
			Owner         string `json:"owner"`
			Assets        []struct {
				Blockchain      string `json:"blockchain"`
				Name            string `json:"name"`
				TokenId         string `json:"tokenId"`
				TokenUrl        string `json:"tokenUrl"`
				ImageUrl        string `json:"imageUrl"`
				CollectionName  string `json:"collectionName"`
				Symbol          string `json:"symbol"`
				ContractType    string `json:"contractType"`
				ContractAddress string `json:"contractAddress"`
				Quantity        string `json:"quantity"`
				Traits          []struct {
					TraitType string      `json:"trait_type"`
					Value     interface{} `json:"value"`
				} `json:"traits"`
			} `json:"assets"`
			SyncStatus interface{} `json:"syncStatus"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Error != nil && response.Error.Message != "" {
		return nil, "", fmt.Errorf("ankr api error: %s", response.Error.Message)
	}

	// 处理NFT数据
	nfts := make([]api.NFT, 0)
	for _, asset := range response.Result.Assets {
		// 只处理ERC721和ERC1155类型的token
		if asset.ContractType != "ERC721" && asset.ContractType != "ERC1155" {
			continue
		}

		fmt.Printf("Found NFT: %s (Type: %s, TokenId: %s)\n", asset.Name, asset.ContractType, asset.TokenId)

		// 处理traits
		attributes := make([]api.Attribute, 0)
		for _, trait := range asset.Traits {
			attributes = append(attributes, api.Attribute{
				TraitType: trait.TraitType,
				Value:     trait.Value,
			})
		}

		nft := api.NFT{
			ContractAddress: asset.ContractAddress,
			TokenId:         asset.TokenId,
			Type:            asset.ContractType,
			Name:            asset.Name,
			Description:     "", // 响应中没有description字段
			Image:           asset.ImageUrl,
			Attributes:      attributes,
			Collection:      asset.CollectionName,
			TokenUri:        asset.TokenUrl,
		}

		nfts = append(nfts, nft)
	}

	fmt.Printf("Returning %d NFTs for address %s\n", len(nfts), address)
	return nfts, response.Result.NextPageToken, nil
}
