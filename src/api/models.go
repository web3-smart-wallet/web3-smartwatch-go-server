package api

// NFT 相关结构
type NFTMetadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Attributes  []Attribute `json:"attributes"`
}

type NFT struct {
	ContractAddress string      `json:"contractAddress"`
	TokenId         string      `json:"tokenId"`
	Type            string      `json:"type"` // ERC721 or ERC1155
	Name            string      `json:"name,omitempty"`
	Description     string      `json:"description,omitempty"`
	Image           string      `json:"image,omitempty"`
	Attributes      []Attribute `json:"attributes,omitempty"`
	Collection      string      `json:"collection,omitempty"`
	TokenUri        string      `json:"tokenUri,omitempty"`
}

type Attribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

type GetApiUserAddressNFTsResponse struct {
	Address    string     `json:"address"`
	NFTs       []NFT      `json:"nfts"`
	Pagination Pagination `json:"pagination"`
}
