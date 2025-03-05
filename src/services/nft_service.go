package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/web3-smart-wallet/src/api"
)

type NFTService struct {
	apiURL string
}

type NFTServiceInterface interface {
	GetNFTs(address string, includeMetadata bool) ([]api.NFT, error)
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

func (s *NFTService) GetNFTs(address string, includeMetadata bool) ([]api.NFT, error) {
	log.Printf("Fetching NFTs for address: %s", address)
	// 使用 eth_getLogs 查询 ERC721/ERC1155 的 Transfer 事件
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getLogs",
		"params": []interface{}{
			map[string]interface{}{
				"fromBlock": "0x1000000",
				"toBlock":   "latest",
				"topics": []interface{}{
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", // ERC721 Transfer
					nil,
					nil,
					"0x000000000000000000000000" + strings.ToLower(strings.TrimPrefix(address, "0x")),
				},
			},
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(s.apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch NFTs: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Result []struct {
			Address string   `json:"address"`
			Topics  []string `json:"topics"`
			Data    string   `json:"data"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("Found %d transfer events", len(response.Result))

	nfts := make([]api.NFT, 0)
	for _, event := range response.Result {
		// 获取 NFT 合约信息
		contractType := s.getContractType(event.Address)
		if contractType == "" {
			continue // 跳过非 NFT 合约
		}

		tokenId := s.getTokenId(event.Data, event.Topics)
		if tokenId == "" {
			continue
		}

		nft := api.NFT{
			TokenAddress: event.Address,
			TokenId:      tokenId,
			Type:         contractType,
		}

		if includeMetadata {
			metadata, err := s.getNFTMetadata(event.Address, tokenId)
			if err != nil {
				log.Printf("Failed to fetch metadata for NFT %s-%s: %v", event.Address, tokenId, err)
				continue
			}
			nft.Metadata = metadata
		}

		nfts = append(nfts, nft)
	}

	return nfts, nil
}

func (s *NFTService) getContractType(address string) string {
	// 先检查 ERC721
	if s.supportsInterface(address, "0x80ac58cd") {
		return "ERC721"
	}
	// 如果不支持 supportsInterface，尝试调用 name() 和 symbol()
	if s.hasERC721Methods(address) {
		return "ERC721"
	}
	return ""
}

func (s *NFTService) supportsInterface(address, interfaceId string) bool {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params": []interface{}{
			map[string]interface{}{
				"to":   address,
				"data": fmt.Sprintf("0x01ffc9a7%s", strings.TrimPrefix(interfaceId, "0x")),
			},
			"latest",
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return false
	}

	resp, err := http.Post(s.apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to check interface: %v", err)
		return false
	}
	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return false
	}

	// 检查返回值是否为 true (0x0000...0001)
	return strings.HasSuffix(response.Result, "01")
}

func (s *NFTService) hasERC721Methods(address string) bool {
	// 调用 name()
	nameData := "0x06fdde03"
	// 调用 symbol()
	symbolData := "0x95d89b41"

	// 如果这两个方法都存在，很可能是 ERC721
	return s.methodExists(address, nameData) && s.methodExists(address, symbolData)
}

func (s *NFTService) methodExists(address, data string) bool {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params": []interface{}{
			map[string]interface{}{
				"to":   address,
				"data": data,
			},
			"latest",
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	resp, err := http.Post(s.apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func (s *NFTService) getTokenId(data string, topics []string) string {
	// 从日志数据中提取 tokenId
	if len(topics) >= 4 {
		// 移除前导零
		return strings.TrimLeft(topics[3], "0")
	}
	return ""
}

func (s *NFTService) getNFTMetadata(address, tokenId string) (*api.NFTMetadata, error) {
	// 获取 tokenURI
	data := fmt.Sprintf("0xc87b56dd%064s", strings.TrimPrefix(tokenId, "0x"))

	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params": []interface{}{
			map[string]interface{}{
				"to":   address,
				"data": data,
			},
			"latest",
		},
		"id": 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(s.apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token URI: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// 解码返回的字符串数据
	uri := decodeString(response.Result)
	if !strings.HasPrefix(uri, "http") {
		uri = "ipfs://" + strings.TrimPrefix(uri, "ipfs://")
	}

	// 获取元数据
	metadataResp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %v", err)
	}
	defer metadataResp.Body.Close()

	var metadata api.NFTMetadata
	if err := json.NewDecoder(metadataResp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %v", err)
	}

	return &metadata, nil
}

// 辅助函数：解码智能合约返回的字符串
func decodeString(hexStr string) string {
	// 移除 0x 前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// 获取字符串长度（前32字节中的最后一个字节）
	offset := new(big.Int)
	offset.SetString(hexStr[:64], 16)

	// 获取字符串内容
	length := new(big.Int)
	length.SetString(hexStr[offset.Int64()*2:offset.Int64()*2+64], 16)

	// 解码字符串内容
	contentStart := offset.Int64()*2 + 64
	contentEnd := contentStart + length.Int64()*2
	content := hexStr[contentStart:contentEnd]

	bytes := make([]byte, length.Int64())
	for i := 0; i < len(bytes); i++ {
		b, _ := strconv.ParseUint(content[i*2:(i+1)*2], 16, 8)
		bytes[i] = byte(b)
	}

	return string(bytes)
}
