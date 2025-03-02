package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/web3-smart-wallet/src/api"
)

type QuickNodeService struct {
	providerURL string
}

type QuickNodeServiceInterface interface {
	GetTokens(address string, includeZeroBalance bool) ([]api.Token, error)
}

type tokenResponse struct {
	Result []struct {
		Address  string `json:"address"`
		Symbol   string `json:"symbol"`
		Name     string `json:"name"`
		Decimals int32  `json:"decimals"`
		Balance  string `json:"balance"`
		Type     string `json:"type"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewQuickNodeService() QuickNodeServiceInterface {
	providerURL := os.Getenv("QUICKNODE_PROVIDER_URL")
	if providerURL == "" {
		panic("QUICKNODE_PROVIDER_URL environment variable is not set")
	}

	return &QuickNodeService{
		providerURL: providerURL,
	}
}

func (s *QuickNodeService) GetTokens(address string, includeZeroBalance bool) ([]api.Token, error) {
	// 构建 JSON-RPC 请求
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "qn_getWalletTokenBalance",
		"params":  []interface{}{address, includeZeroBalance},
		"id":      1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 发送请求
	resp, err := http.Post(s.providerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Error != nil && response.Error.Message != "" {
		return nil, fmt.Errorf("quicknode error: %s", response.Error.Message)
	}

	// 转换为 API Token 格式
	tokens := make([]api.Token, 0)
	for _, t := range response.Result {
		tokenType := api.ERC20
		if t.Type != "" {
			tokenType = api.TokenType(t.Type)
		}

		token := api.Token{
			Address:  t.Address,
			Symbol:   t.Symbol,
			Name:     t.Name,
			Decimals: int(t.Decimals),
			Type:     &tokenType,
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}
