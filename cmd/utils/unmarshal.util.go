package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/unicornultrafoundation/go-u2u/common"

	"github.com/google/uuid"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

// Define a custom type
type JsonNullInt16 struct {
	sql.NullInt16
}

// Implement the UnmarshalJSON method
func (v *JsonNullInt16) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal into a float64 first
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int16 = int16(*x)
	} else {
		v.Valid = false
	}
	return nil
}

// Define a custom type
type JsonNullTime struct {
	sql.NullTime
}

// Implement the UnmarshalJSON method
func (v *JsonNullTime) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal into a string
	var str *string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != nil {
		// Parse the time using the desired layout
		t, err := time.Parse(time.DateTime, *str)
		if err != nil {
			return err
		}
		v.Valid = true
		v.Time = t
	} else {
		v.Valid = false
	}

	return nil
}

type AddNewAssetParamsUtil struct {
	ID                string
	ChainID           int32
	CollectionAddress string
	Type              db.AssetType
	DecimalData       JsonNullInt16
	InitialBlock      JsonNullInt64
	LastUpdated       JsonNullTime
}

type AddNewAssetParamsSwagger struct {
	ChainID           int32
	CollectionAddress string
	Type              db.AssetType
	DecimalData       int16
	InitialBlock      int64
	LastUpdated       time.Time
}

// AddNewAssetResponse represents the response parameters for adding a new asset
type AssetResponse struct {
	ID                string       `json:"id"`
	ChainID           int32        `json:"chainId"`
	CollectionAddress string       `json:"collectionAddress"`
	Type              db.AssetType `json:"type"`
	DecimalData       int          `json:"decimalData"`  // Output as int
	InitialBlock      int64        `json:"initialBlock"` // Output as int64
	LastUpdated       time.Time    `json:"lastUpdated"`  // Output as time.Time
}

// Convert function
func ConvertCustomTypeToSqlParams(param *AddNewAssetParamsUtil) db.AddNewAssetParams {
	return db.AddNewAssetParams{
		ID:                param.ID,
		ChainID:           param.ChainID,
		CollectionAddress: common.HexToAddress(param.CollectionAddress).Hex(),
		Type:              param.Type,
		DecimalData:       sql.NullInt16{Int16: param.DecimalData.Int16, Valid: param.DecimalData.Valid},
		InitialBlock:      sql.NullInt64{Int64: param.InitialBlock.Int64, Valid: param.InitialBlock.Valid},
		LastUpdated:       sql.NullTime{Time: param.LastUpdated.Time, Valid: param.LastUpdated.Valid},
	}
}

// MarshalParams marshals the AddNewAssetParams struct to a JSON object
func MarshalAssetParams(params db.AddNewAssetParams) (AssetResponse, error) {
	// Prepare the response struct
	response := AssetResponse{
		ID:                params.ID,
		ChainID:           params.ChainID,
		CollectionAddress: params.CollectionAddress,
		Type:              params.Type,
		LastUpdated:       params.LastUpdated.Time,
	}

	// Convert nullable DecimalData
	if params.DecimalData.Valid {
		response.DecimalData = int(params.DecimalData.Int16) // Convert to int
	} else {
		response.DecimalData = 0 // Default value if null
	}

	// Convert nullable InitialBlock
	if params.InitialBlock.Valid {
		response.InitialBlock = params.InitialBlock.Int64 // Convert to int64
	} else {
		response.InitialBlock = 0 // Default value if null
	}

	return response, nil // Return the response struct
}

func ConvertToAssetResponses(assets []db.Asset) []AssetResponse {
	responses := make([]AssetResponse, 0, len(assets))

	for _, asset := range assets {
		response := AssetResponse{
			ID:                asset.ID,
			ChainID:           asset.ChainID,
			CollectionAddress: asset.CollectionAddress,
			Type:              asset.Type,
		}

		// Handle DecimalData conversion
		if asset.DecimalData.Valid {
			response.DecimalData = int(asset.DecimalData.Int16)
		}

		// Handle InitialBlock conversion
		if asset.InitialBlock.Valid {
			response.InitialBlock = asset.InitialBlock.Int64
		}

		// Handle LastUpdated conversion
		if asset.LastUpdated.Valid {
			response.LastUpdated = asset.LastUpdated.Time
		}

		responses = append(responses, response)
	}

	return responses
}

func ConvertAssetToAssetResponse(asset db.Asset) AssetResponse {
	response := AssetResponse{
		ID:                asset.ID,
		ChainID:           asset.ChainID,
		CollectionAddress: asset.CollectionAddress,
		Type:              asset.Type,
	}

	// Handle DecimalData conversion
	if asset.DecimalData.Valid {
		response.DecimalData = int(asset.DecimalData.Int16)
	}

	// Handle InitialBlock conversion
	if asset.InitialBlock.Valid {
		response.InitialBlock = asset.InitialBlock.Int64
	}

	// Handle LastUpdated conversion
	if asset.LastUpdated.Valid {
		response.LastUpdated = asset.LastUpdated.Time
	}

	return response
}

type Erc721CollectionAssetResponse struct {
	ID         uuid.UUID `json:"id"`
	ChainID    int32     `json:"chainId"`
	AssetID    string    `json:"assetId"`
	TokenID    string    `json:"tokenId"`
	Owner      string    `json:"owner"`
	Attributes string    `json:"attributes"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func ConvertToErc721CollectionAssetResponses(assets []db.Erc721CollectionAsset) []Erc721CollectionAssetResponse {
	responses := make([]Erc721CollectionAssetResponse, 0, len(assets))

	for _, asset := range assets {
		response := Erc721CollectionAssetResponse{
			ID:         asset.ID,
			ChainID:    asset.ChainID,
			AssetID:    asset.AssetID,
			TokenID:    asset.TokenID,
			Owner:      asset.Owner,
			Attributes: asset.Attributes.String,
			CreatedAt:  asset.CreatedAt,
			UpdatedAt:  asset.UpdatedAt,
		}

		responses = append(responses, response)
	}

	return responses
}

func ConvertToDetailErc721CollectionAssetResponses(asset db.Erc721CollectionAsset) Erc721CollectionAssetResponse {
	response := Erc721CollectionAssetResponse{
		ID:         asset.ID,
		ChainID:    asset.ChainID,
		AssetID:    asset.AssetID,
		TokenID:    asset.TokenID,
		Owner:      asset.Owner,
		Attributes: asset.Attributes.String,
		CreatedAt:  asset.CreatedAt,
		UpdatedAt:  asset.UpdatedAt,
	}

	return response
}

type Erc1155CollectionAssetResponse struct {
	ID         uuid.UUID `json:"id"`
	ChainID    int32     `json:"chainId"`
	AssetID    string    `json:"assetId"`
	TokenID    string    `json:"tokenId"`
	Owner      string    `json:"owner"`
	Balance    string    `json:"balance"`
	Attributes string    `json:"attributes"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func ConvertToErc1155CollectionAssetResponses(assets []db.Erc1155CollectionAsset) []Erc1155CollectionAssetResponse {
	responses := make([]Erc1155CollectionAssetResponse, 0, len(assets))

	for _, asset := range assets {
		response := Erc1155CollectionAssetResponse{
			ID:         asset.ID,
			ChainID:    asset.ChainID,
			AssetID:    asset.AssetID,
			TokenID:    asset.TokenID,
			Owner:      asset.Owner,
			Balance:    asset.Balance,
			Attributes: asset.Attributes.String,
			CreatedAt:  asset.CreatedAt,
			UpdatedAt:  asset.UpdatedAt,
		}

		responses = append(responses, response)
	}

	return responses
}

// CustomTime is a wrapper around time.Time to handle custom unmarshalling.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface for CustomTime.
func (c *CustomTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// Parse the time string according to the expected layout
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	c.Time = t
	return nil
}

type ERC1155AssetOwner struct {
	Balance   int        `json:"balance"`
	CreatedAt CustomTime `json:"createdAt"`
	Id        uuid.UUID  `json:"id"`
	Owner     string     `json:"owner"`
	UpdatedAt CustomTime `json:"updatedAt"`
}

type ERC1155AssetOwnerResponse struct {
	Balance   string     `json:"balance"`
	CreatedAt CustomTime `json:"createdAt"`
	Id        uuid.UUID  `json:"id"`
	Owner     string     `json:"owner"`
	UpdatedAt CustomTime `json:"updatedAt"`
}

type GetDetailERC1155Asset struct {
	AssetID     string                      `json:"assetId"`
	TokenID     string                      `json:"tokenId"`
	Attributes  string                      `json:"attributes"`
	TotalSupply string                      `json:"totalSupply"`
	AssetOwners []ERC1155AssetOwnerResponse `json:"assetOwners"`
}

func ConvertERC1155Owner(rawData []byte) ([]ERC1155AssetOwnerResponse, error) {
	var owners []ERC1155AssetOwner
	var results []ERC1155AssetOwnerResponse

	// Unmarshal the JSON array directly into the slice
	err := json.Unmarshal(rawData, &owners)

	for _, owner := range owners {
		results = append(results, ERC1155AssetOwnerResponse{
			Balance:   strconv.Itoa(owner.Balance),
			CreatedAt: CustomTime{owner.CreatedAt.Time},
			Id:        owner.Id,
			Owner:     owner.Owner,
			UpdatedAt: CustomTime{owner.UpdatedAt.Time},
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ERC1155 owners: %w", err)
	}

	return results, nil
}

func ConvertToDetailERC1155AssetResponse(asset db.GetDetailERC1155AssetsRow) GetDetailERC1155Asset {

	response := GetDetailERC1155Asset{
		AssetID:     asset.AssetID,
		TokenID:     asset.TokenID,
		Attributes:  asset.Attributes.String,
		TotalSupply: strconv.FormatInt(asset.TotalSupply, 10),
		AssetOwners: func() []ERC1155AssetOwnerResponse {
			owners, err := ConvertERC1155Owner(asset.AssetOwners)
			if err != nil {
				return nil
			}
			return owners
		}(),
	}

	return response
}
