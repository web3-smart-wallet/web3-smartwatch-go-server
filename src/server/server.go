package server

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/web3-smart-wallet/src/api"
	"github.com/web3-smart-wallet/src/services"
)

// 添加地址验证正则表达式
var addressRegex = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")

type Server struct {
	ankrService services.AnkrServiceInterface
	nftService  services.NFTServiceInterface
}

func NewServer(ankrService services.AnkrServiceInterface, nftService services.NFTServiceInterface) api.ServerInterface {
	return &Server{
		ankrService: ankrService,
		nftService:  nftService,
	}
}

func (s Server) GetApiSearchAddressAddress(c *fiber.Ctx, address string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetApiSearchDidDid(c *fiber.Ctx, did string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetApiUserAddress(c *fiber.Ctx, address string, params api.GetApiUserAddressParams) error {
	// 验证地址格式
	if !addressRegex.MatchString(address) {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    "invalid_address",
			Message: "Invalid Ethereum address format",
		})
	}

	// 调用服务获取代币列表
	tokens, err := s.ankrService.GetTokenList(address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    "internal_server_error",
			Message: err.Error(),
		})
	}

	// 返回响应
	return c.JSON(fiber.Map{
		"address": address,
		"tokens":  tokens,
	})
}

func (s Server) GetApiUserAddressNfts(c *fiber.Ctx, address string, params api.GetApiUserAddressNftsParams) error {
	// 验证地址格式
	if !addressRegex.MatchString(address) {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    "invalid_address",
			Message: "Invalid Ethereum address format",
		})
	}

	fmt.Printf("Fetching NFTs for address: %s\n", address)

	// 获取NFT列表
	includeMetadata := true
	if params.IncludeMetadata != nil {
		includeMetadata = *params.IncludeMetadata
	}

	fmt.Printf("Include metadata: %v\n", includeMetadata)

	nfts, err := s.nftService.GetNFTs(address, includeMetadata)
	if err != nil {
		fmt.Printf("Error fetching NFTs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    "internal_server_error",
			Message: err.Error(),
		})
	}

	fmt.Printf("Found %d NFTs for address %s\n", len(nfts), address)

	// 返回响应
	return c.JSON(fiber.Map{
		"address": address,
		"nfts":    nfts,
	})
}

func (s Server) GetApiUserAddressBalance(c *fiber.Ctx, address string, params api.GetApiUserAddressBalanceParams) error {
	// 验证地址格式
	if !addressRegex.MatchString(address) {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    "invalid_address",
			Message: "Invalid Ethereum address format",
		})
	}

	// 获取代币余额
	includeZeroBalance := c.Query("includeZeroBalance") == "true"

	// 调用服务获取代币信息
	tokens, err := s.ankrService.GetTokens(address, includeZeroBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    "internal_server_error",
			Message: err.Error(),
		})
	}

	// 返回响应
	return c.JSON(fiber.Map{
		"address": address,
		"tokens":  tokens, // 直接返回tokens，不再转换为balances
	})
}

// 辅助函数：格式化代币余额
func formatBalance(hexBalance string, decimals int) string {
	// 移除 0x 前缀
	hexBalance = strings.TrimPrefix(hexBalance, "0x")

	// 转换为大整数
	balance := new(big.Int)
	balance.SetString(hexBalance, 16)

	// 创建 10^decimals 的除数
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// 计算整数部分和小数部分
	intPart := new(big.Int).Div(balance, divisor)
	remainder := new(big.Int).Mod(balance, divisor)

	// 格式化小数部分
	decimalStr := remainder.String()
	// 补齐前导零
	for len(decimalStr) < decimals {
		decimalStr = "0" + decimalStr
	}

	// 移除尾部的零
	decimalStr = strings.TrimRight(decimalStr, "0")

	// 组合结果
	if decimalStr == "" {
		return intPart.String()
	}
	return intPart.String() + "." + decimalStr
}
