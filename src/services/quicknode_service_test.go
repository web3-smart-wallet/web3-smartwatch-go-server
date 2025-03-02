package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuickNodeService_GetTokens(t *testing.T) {
	service := NewQuickNodeService()

	// 测试用例1：获取非零余额代币
	tokens, err := service.GetTokens("0x742d35Cc6634C0532925a3b844Bc454e4438f44e", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens)

	// 测试用例2：包含零余额代币
	tokensWithZero, err := service.GetTokens("0x742d35Cc6634C0532925a3b844Bc454e4438f44e", true)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokensWithZero)

	// 验证返回的代币格式
	for _, token := range tokens {
		assert.NotEmpty(t, token.Address)
		assert.NotEmpty(t, token.Symbol)
		assert.NotEmpty(t, token.Name)
		assert.NotNil(t, token.Type)
	}
}
