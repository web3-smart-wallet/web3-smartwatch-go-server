package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/web3-smart-wallet/src/api"
)

type AnkrService struct {
	apiURL string
}

type AnkrServiceInterface interface {
	GetTokens(address string, includeZeroBalance bool) ([]api.Token, error)
	GetTokenList(address string) ([]api.Token, error)
}

func NewAnkrService(apiURL string) AnkrServiceInterface {
	return &AnkrService{
		apiURL: apiURL,
	}
}

// 定义已知的代币列表
var knownTokens = []struct {
	Address  string
	Symbol   string
	Name     string
	Decimals int
}{
	{
		Address:  "0x4200000000000000000000000000000000000006",
		Symbol:   "WETH",
		Name:     "Wrapped Ether",
		Decimals: 18,
	},
	{
		Address:  "0xd9aAEc86B65D86f6A7B5B1b0c42FFA531710b6CA",
		Symbol:   "USDbC",
		Name:     "USD Base Coin",
		Decimals: 6,
	},
}

func (s *AnkrService) GetTokens(address string, includeZeroBalance bool) ([]api.Token, error) {
	tokens := make([]api.Token, 0)

	// 遍历已知代币列表
	for _, token := range knownTokens {
		// 构建 balanceOf 调用数据
		data := fmt.Sprintf("0x70a08231000000000000000000000000%s", strings.TrimPrefix(address, "0x"))

		payload := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_call",
			"params": []interface{}{
				map[string]interface{}{
					"to":   token.Address,
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
			return nil, fmt.Errorf("failed to fetch token balance: %v", err)
		}
		defer resp.Body.Close()

		var response struct {
			Result string `json:"result"`
			Error  *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %v", err)
		}

		if response.Error != nil && response.Error.Message != "" {
			return nil, fmt.Errorf("rpc error: %s", response.Error.Message)
		}

		// 如果余额不为0或includeZeroBalance为true
		if includeZeroBalance || response.Result != "0x0" {
			tokenType := api.TokenType("ERC20")
			balance := response.Result
			tokens = append(tokens, api.Token{
				Address:  token.Address,
				Symbol:   token.Symbol,
				Name:     token.Name,
				Decimals: token.Decimals,
				Type:     &tokenType,
				Balance:  &balance,
			})
		}
	}

	log.Printf("Found %d tokens", len(tokens))
	return tokens, nil
}

func (s *AnkrService) GetTokenList(address string) ([]api.Token, error) {
	// 构建请求体
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "ankr_getAccountBalance",
		"params": map[string]interface{}{
			"blockchain":      "base",
			"walletAddress":   address,
			"onlyWhitelisted": false,
		},
		"id": 1,
	}

	// 发送请求
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response struct {
		Result struct {
			Assets []struct {
				TokenName   string `json:"tokenName"`
				TokenSymbol string `json:"tokenSymbol"`
				Address     string `json:"contractAddress"`
				Decimals    int    `json:"tokenDecimals"`
			} `json:"assets"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Error != nil && response.Error.Message != "" {
		return nil, fmt.Errorf("ankr api error: %s", response.Error.Message)
	}

	// 转换为 API 响应格式
	tokens := make([]api.Token, len(response.Result.Assets))
	for i, asset := range response.Result.Assets {
		tokens[i] = api.Token{
			Address:  asset.Address,
			Symbol:   asset.TokenSymbol,
			Name:     asset.TokenName,
			Decimals: asset.Decimals,
		}
	}

	return tokens, nil
}
