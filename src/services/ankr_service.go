package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/web3-smart-wallet/src/api"
)

type AnkrService struct {
	apiURL string
}

type AnkrServiceInterface interface {
	GetTokens(address string, includeZeroBalance bool) ([]api.Token, error)
}

func NewAnkrService() AnkrServiceInterface {
	apiURL := os.Getenv("ANKR_API_URL")
	if apiURL == "" {
		panic("ANKR_API_URL environment variable is not set")
	}

	return &AnkrService{
		apiURL: apiURL,
	}
}

func (s *AnkrService) GetTokens(address string, includeZeroBalance bool) ([]api.Token, error) {
	// 定义要查询的代币列表
	tokenConfigs := []struct {
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

	tokens := make([]api.Token, 0)
	tokenType := api.TokenType("ERC20")

	// 遍历每个代币合约查询余额
	for _, config := range tokenConfigs {
		// 构建 balanceOf 调用数据
		data := "0x70a08231000000000000000000000000" + address[2:]

		payload := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_call",
			"params": []interface{}{
				map[string]interface{}{
					"to":   config.Address, // 调用代币合约
					"data": data,           // balanceOf(address)
				},
				"latest",
			},
			"id": 1,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %v", err)
		}

		log.Printf("Fetching balance for token %s at address %s", config.Symbol, config.Address)
		log.Printf("Request payload: %+v", payload)

		// 发送请求
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
			return nil, fmt.Errorf("ankr error: %s", response.Error.Message)
		}

		log.Printf("Response for %s: %s", config.Symbol, response.Result)

		// 如果余额不为0或includeZeroBalance为true，添加到结果中
		balance := response.Result
		if includeZeroBalance || balance != "0x0" {
			tokens = append(tokens, api.Token{
				Address:  config.Address,
				Symbol:   config.Symbol,
				Name:     config.Name,
				Decimals: config.Decimals,
				Type:     &tokenType,
				Balance:  &balance,
			})
		}
	}

	log.Printf("Found tokens: %+v", tokens)
	return tokens, nil
}
