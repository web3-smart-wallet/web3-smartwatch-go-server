package server

import (
	"regexp"

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
	//TODO implement me
	panic("implement me")
}

func (s Server) GetApiUserAddressBalance(c *fiber.Ctx, address string, params api.GetApiUserAddressBalanceParams) error {
	//TODO implement me
	panic("implement me")
}
