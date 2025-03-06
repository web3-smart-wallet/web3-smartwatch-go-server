package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

func (s *AnkrService) GetTokens(address string, includeZeroBalance bool) ([]api.Token, error) {
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
				Balance     string `json:"balance"`
				BalanceUsd  string `json:"balanceUsd,omitempty"`
				TokenPrice  string `json:"tokenPrice,omitempty"`
				TokenType   string `json:"tokenType"`
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
	tokens := make([]api.Token, 0)
	for _, asset := range response.Result.Assets {
		// 如果余额为0且不包含零余额，则跳过
		if !includeZeroBalance && (asset.Balance == "0" || asset.Balance == "0x0") {
			continue
		}

		tokenType := api.TokenType(asset.TokenType)
		if tokenType == "" {
			tokenType = api.TokenType("ERC20")
		}

		// 创建Token对象，包含所需字段
		token := api.Token{
			Address:  asset.Address,
			Symbol:   asset.TokenSymbol,
			Name:     asset.TokenName,
			Decimals: asset.Decimals,
			Type:     &tokenType,
			Balance:  &asset.Balance,
		}

		tokens = append(tokens, token)
	}

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

	// 解析响应 - 更新结构体以包含更多字段
	var response struct {
		Result struct {
			Assets []struct {
				TokenName   string `json:"tokenName"`
				TokenSymbol string `json:"tokenSymbol"`
				Address     string `json:"contractAddress"`
				Decimals    int    `json:"tokenDecimals"`
				Balance     string `json:"balance,omitempty"`
				BalanceUsd  string `json:"balanceUsd,omitempty"`
				TokenPrice  string `json:"tokenPrice,omitempty"`
				TokenType   string `json:"tokenType,omitempty"`
				Thumbnail   string `json:"thumbnail,omitempty"`
				IsVerified  bool   `json:"isVerified,omitempty"`
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
		tokenType := api.TokenType(asset.TokenType)
		if tokenType == "" {
			tokenType = api.TokenType("ERC20")
		}

		token := api.Token{
			Address:  asset.Address,
			Symbol:   asset.TokenSymbol,
			Name:     asset.TokenName,
			Decimals: asset.Decimals,
			Type:     &tokenType,
		}

		// 如果有余额信息，添加到Token中
		if asset.Balance != "" {
			token.Balance = &asset.Balance
		}

		tokens[i] = token
	}

	return tokens, nil
}
