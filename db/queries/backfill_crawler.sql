-- name: GetCrawlingBackfillCrawler :many
SELECT * FROM backfill_crawlers 
WHERE status = crawler_status('CRAWLING');

-- name: UpdateCrawlingBackfill :exec
UPDATE backfill_crawlers
SET 
    status = COALESCE($2, status),            
    current_block = COALESCE($3, current_block)  
WHERE id = $1;

-- name: AddBackfillCrawler :exec
INSERT INTO backfill_crawlers (
    chain_id, collection_address, current_block
)
VALUES (
    $1, $2, $3
) RETURNING *;
