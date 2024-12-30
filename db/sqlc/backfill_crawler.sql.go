// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: backfill_crawler.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const addBackfillCrawler = `-- name: AddBackfillCrawler :exec
INSERT INTO backfill_crawlers (
    chain_id, collection_address, current_block
)
VALUES (
    $1, $2, $3
) ON CONFLICT ON CONSTRAINT BACKFILL_CRAWLERS_PKEY DO UPDATE SET
    current_block = EXCLUDED.current_block,
    status = 'CRAWLING'
RETURNING chain_id, collection_address, current_block, status, created_at
`

type AddBackfillCrawlerParams struct {
	ChainID           int32  `json:"chainId"`
	CollectionAddress string `json:"collectionAddress"`
	CurrentBlock      int64  `json:"currentBlock"`
}

func (q *Queries) AddBackfillCrawler(ctx context.Context, arg AddBackfillCrawlerParams) error {
	_, err := q.db.ExecContext(ctx, addBackfillCrawler, arg.ChainID, arg.CollectionAddress, arg.CurrentBlock)
	return err
}

const getCrawlingBackfillCrawler = `-- name: GetCrawlingBackfillCrawler :many
SELECT 
    bc.chain_id, bc.collection_address, bc.current_block, bc.status, bc.created_at, 
    a.type, 
    a.initial_block 
FROM 
    backfill_crawlers AS bc
JOIN 
    assets AS a 
    ON a.chain_id = bc.chain_id 
    AND a.collection_address = bc.collection_address 
WHERE 
    bc.status = 'CRAWLING'
`

type GetCrawlingBackfillCrawlerRow struct {
	ChainID           int32         `json:"chainId"`
	CollectionAddress string        `json:"collectionAddress"`
	CurrentBlock      int64         `json:"currentBlock"`
	Status            CrawlerStatus `json:"status"`
	CreatedAt         time.Time     `json:"createdAt"`
	Type              AssetType     `json:"type"`
	InitialBlock      sql.NullInt64 `json:"initialBlock"`
}

func (q *Queries) GetCrawlingBackfillCrawler(ctx context.Context) ([]GetCrawlingBackfillCrawlerRow, error) {
	rows, err := q.db.QueryContext(ctx, getCrawlingBackfillCrawler)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCrawlingBackfillCrawlerRow
	for rows.Next() {
		var i GetCrawlingBackfillCrawlerRow
		if err := rows.Scan(
			&i.ChainID,
			&i.CollectionAddress,
			&i.CurrentBlock,
			&i.Status,
			&i.CreatedAt,
			&i.Type,
			&i.InitialBlock,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCrawlingBackfill = `-- name: UpdateCrawlingBackfill :exec
UPDATE backfill_crawlers
SET 
    status = COALESCE($3, status),            
    current_block = COALESCE($4, current_block)  
WHERE chain_id = $1
AND collection_address = $2
`

type UpdateCrawlingBackfillParams struct {
	ChainID           int32         `json:"chainId"`
	CollectionAddress string        `json:"collectionAddress"`
	Status            CrawlerStatus `json:"status"`
	CurrentBlock      int64         `json:"currentBlock"`
}

func (q *Queries) UpdateCrawlingBackfill(ctx context.Context, arg UpdateCrawlingBackfillParams) error {
	_, err := q.db.ExecContext(ctx, updateCrawlingBackfill,
		arg.ChainID,
		arg.CollectionAddress,
		arg.Status,
		arg.CurrentBlock,
	)
	return err
}
