package server

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/web3-smart-wallet/src/api"
	"github.com/web3-smart-wallet/src/services"
	"github.com/web3-smart-wallet/src/utils"
)

// 添加地址验证正则表达式
var addressRegex = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")

type Server struct {
	ankrService services.AnkrServiceInterface
}

func NewServer(ankrService services.AnkrServiceInterface) api.ServerInterface {
	return &Server{
		ankrService: ankrService,
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

	tokens, err := s.ankrService.GetTokens(address, params.IncludeZeroBalance != nil && *params.IncludeZeroBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    "internal_server_error",
			Message: err.Error(),
		})
	}

	// 转换为 TokenBalance 数组
	balances := make([]api.TokenBalance, len(tokens))
	for i, token := range tokens {
		var rawBalance string
		var formattedBalance string
		if token.Balance != nil {
			rawBalance = *token.Balance
			formattedBalance = utils.FormatTokenBalance(rawBalance, token.Decimals)
		}
		balances[i] = api.TokenBalance{
			Token:            token,
			RawBalance:       rawBalance,
			FormattedBalance: &formattedBalance,
		}
	}

	return c.JSON(fiber.Map{
		"address":  address,
		"balances": balances,
		"pagination": map[string]interface{}{
			"current_page":   1,
			"total_pages":    1,
			"total_items":    len(balances),
			"items_per_page": len(balances),
		},
	})
}

func (s Server) GetApiUserAddressNFTs(c *fiber.Ctx, address string, params api.GetApiUserAddressNFTsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetApiUserAddressBalance(c *fiber.Ctx, address string, params api.GetApiUserAddressBalanceParams) error {
	//TODO implement me
	panic("implement me")
}
