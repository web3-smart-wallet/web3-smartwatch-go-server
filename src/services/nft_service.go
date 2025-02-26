package services

type NftService struct {
}

type NftServiceInterface interface {
}

func NewNftService() NftServiceInterface {
	return &NftService{}
}
