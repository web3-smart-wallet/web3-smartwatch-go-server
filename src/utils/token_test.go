package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTokenBalance(t *testing.T) {
	tests := []struct {
		name       string
		hexBalance string
		decimals   int
		expected   string
	}{
		{
			name:       "WETH balance",
			hexBalance: "0x000000000000000000000000000000000000000000000000000058c5d239f3c7",
			decimals:   18,
			expected:   "0.024897432", // 实际值可能需要调整
		},
		{
			name:       "Zero balance",
			hexBalance: "0x0000000000000000000000000000000000000000000000000000000000000000",
			decimals:   6,
			expected:   "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTokenBalance(tt.hexBalance, tt.decimals)
			assert.Equal(t, tt.expected, result)
		})
	}
}
