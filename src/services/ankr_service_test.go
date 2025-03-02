package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnkrService_GetTokens(t *testing.T) {
	service := NewAnkrService()

	// 测试用例1：已知钱包地址
	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	tokens, err := service.GetTokens(address, true)
	assert.NoError(t, err)

	// 验证基本数据格式
	for _, token := range tokens {
		// 验证地址格式
		assert.Regexp(t, "^0x[a-fA-F0-9]{40}$", token.Address)
		// 验证必填字段
		assert.NotEmpty(t, token.Symbol)
		assert.NotEmpty(t, token.Name)
		assert.NotNil(t, token.Type)
	}

	// 验证已知代币
	hasWETH := false
	hasUSDbC := false
	for _, token := range tokens {
		if token.Address == "0x4200000000000000000000000000000000000006" {
			hasWETH = true
			assert.Equal(t, "WETH", token.Symbol)
			assert.Equal(t, 18, token.Decimals)
		}
		if token.Address == "0xd9aAEc86B65D86f6A7B5B1b0c42FFA531710b6CA" {
			hasUSDbC = true
			assert.Equal(t, "USDbC", token.Symbol)
			assert.Equal(t, 6, token.Decimals)
		}
	}
	assert.True(t, hasWETH, "WETH token not found")
	assert.True(t, hasUSDbC, "USDbC token not found")
}
