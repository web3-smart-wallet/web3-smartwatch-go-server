package api

// NFT 相关结构
type NFTMetadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Attributes  []Attribute `json:"attributes"`
}

type CustomNFT struct {
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

type CustomPagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

type GetApiUserAddressNFTsResponse struct {
	Address    string           `json:"address"`
	NFTs       []CustomNFT      `json:"nfts"`
	Pagination CustomPagination `json:"pagination"`
}
