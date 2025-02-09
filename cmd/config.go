package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/internal/generator"
	"gopkg.in/yaml.v2"
)

func loadCrawlerConfig() (*generator.CrawlerConfig, error) {
	data, err := os.ReadFile("subgraph.yaml")
	if err != nil {
		return nil, err
	}

	var config generator.CrawlerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// InitializeSystemFromConfig loads the subgraph config into the system database
func InitializeSystemFromConfig(ctx context.Context, queries *db.Queries) error {
	config, err := loadCrawlerConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Convert chainId from string to int64
	chainID, err := strconv.ParseInt(config.Network.ChainId, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid chain ID: %w", err)
	}

	// Extract chain name from repository URL
	chainName := config.Network.ChainId

	// Create or update chain record
	chain, err := queries.CreateChain(ctx, db.CreateChainParams{
		ID:          int32(chainID),
		Chain:       chainName,
		Name:        config.Network.ChainId,
		RpcUrl:      config.Network.Endpoint[0],
		ChainID:     chainID,
		Explorer:    config.Network.Endpoint[0], // Using RPC URL as explorer URL if not specified
		LatestBlock: int64(config.DataSources[0].StartBlock),
		BlockTime:   2, // Default block time
	})

	if err != nil {
		return fmt.Errorf("failed to create chain: %w", err)
	}

	// For each data source, create asset records
	for _, ds := range config.DataSources {
		_, err = queries.CreateAsset(ctx, db.CreateAssetParams{
			ID:              ds.Options.Address,
			ChainID:         chain.ID,
			ContractAddress: ds.Options.Address,
			InitialBlock:    sql.NullInt64{Int64: int64(ds.StartBlock), Valid: true},
		})

		if err != nil {
			return fmt.Errorf("failed to create asset for address %s: %w", ds.Options.Address, err)
		}
	}

	return nil
}
