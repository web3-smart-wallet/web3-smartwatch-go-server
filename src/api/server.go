package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/web3-smart-wallet/src/services"
)

type Server struct {
	nftService services.NftServiceInterface
}

func NewServer(nftService services.NftServiceInterface) ServerInterface {
	return &Server{
		nftService: nftService,
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

func (s Server) GetApiUserAddress(c *fiber.Ctx, address string, params GetApiUserAddressParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetApiUserAddressNFTs(c *fiber.Ctx, address string, params GetApiUserAddressNFTsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetApiUserAddressBalance(c *fiber.Ctx, address string, params GetApiUserAddressBalanceParams) error {
	//TODO implement me
	panic("implement me")
}
