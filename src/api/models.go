package api

// NFT 相关结构
type NFTMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Attributes  []struct {
		TraitType string      `json:"trait_type"`
		Value     interface{} `json:"value"`
	} `json:"attributes,omitempty"`
}

type NFT struct {
	TokenAddress string       `json:"token_address"`
	TokenId      string       `json:"token_id"`
	Type         string       `json:"type"` // ERC721 or ERC1155
	Metadata     *NFTMetadata `json:"metadata,omitempty"`
}

type GetApiUserAddressNFTsResponse struct {
	Address    string     `json:"address"`
	NFTs       []NFT      `json:"nfts"`
	Pagination Pagination `json:"pagination"`
}
