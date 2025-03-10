package utils

import (
	"math/big"
	"strings"
)

// FormatTokenBalance 将十六进制余额转换为可读格式的字符串
// hexBalance: 十六进制格式的余额 (0x...)
// decimals: 代币的小数位数
func FormatTokenBalance(hexBalance string, decimals int) string {
	// 移除 "0x" 前缀
	hex := strings.TrimPrefix(hexBalance, "0x")

	// 转换为大整数
	balance := new(big.Int)
	balance.SetString(hex, 16)

	// 创建 10^decimals 作为除数
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// 计算整数部分和小数部分
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.DivMod(balance, divisor, remainder)

	// 格式化小数部分
	decimalStr := remainder.String()
	// 补齐前导零
	for len(decimalStr) < decimals {
		decimalStr = "0" + decimalStr
	}

	// 如果小数部分全是0，则不显示
	if strings.Trim(decimalStr, "0") == "" {
		return quotient.String()
	}

	// 去掉尾部的0
	decimalStr = strings.TrimRight(decimalStr, "0")

	return quotient.String() + "." + decimalStr
}
