package types

type CombinedAsset struct {
	TokenType  string `json:"token_type"`
	ChainID    int    `json:"chain_id"`
	AssetID    string `json:"asset_id"`
	TokenID    string `json:"token_id"`
	Attributes string `json:"attributes"`
	CreatedAt  string `json:"created_at"`
}
