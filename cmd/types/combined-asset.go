package types

type CombinedAsset struct {
	TokenType  string  `json:"tokenType"`
	ChainID    int     `json:"chainId"`
	AssetID    string  `json:"assetId"`
	TokenID    string  `json:"tokenId"`
	Owner      *string `json:"owner,omitempty"`
	Attributes string  `json:"attributes"`
	CreatedAt  string  `json:"createdAt"`
}
