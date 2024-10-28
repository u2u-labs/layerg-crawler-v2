package db

import (
	"database/sql"
	"strconv"

	"github.com/Masterminds/squirrel"

	"github.com/gin-gonic/gin"
)

// Pagination holds the pagination information.
type Pagination[T any] struct {
	Page       int   `json:"page"`       // Current page number
	Limit      int   `json:"limit"`      // Number of items per page
	TotalItems int64 `json:"totalItems"` // Total number of items available
	TotalPages int64 `json:"totalPages"` // Total number of pages
	Data       []T   `json:"data"`       // The paginated items (can be any type)
}

// get limit and offset from query parameters
func GetLimitAndOffset(ctx *gin.Context) (int, int, int) {
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	// Default values
	pageNum := 1
	limitNum := 10

	// Parse page and limit
	if p, err := strconv.Atoi(pageStr); err == nil {
		pageNum = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil {
		limitNum = l
	}

	// Calculate offset
	offset := (pageNum - 1) * limitNum

	return pageNum, limitNum, offset
}

func FilterItems[T any](db *sql.DB, tableName string, limit int, offset int, filterConditions map[string]string) ([]T, error) {
	// Create a Squirrel query builder for the items table
	queryBuilder := squirrel.Select("*").From(tableName)

	// Apply dynamic filters
	for column, value := range filterConditions {
		queryBuilder = queryBuilder.Where(squirrel.Eq{column: value})
	}

	// Apply pagination
	queryBuilder = queryBuilder.Limit(uint64(limit)).Offset(uint64(offset))

	// Convert the query to SQL
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		var item T
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// CountItems counts the number of items in the database based on dynamic filters.
func CountItemsWithFilter(db *sql.DB, tableName string, filterConditions map[string]string) (int, error) {
	// Create a Squirrel query builder for counting items
	queryBuilder := squirrel.Select("COUNT(*)").From(tableName)

	// Apply dynamic filters
	for column, value := range filterConditions {
		queryBuilder = queryBuilder.Where(squirrel.Eq{column: value})
	}

	// Convert the query to SQL
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	// Execute the query
	var count int
	err = db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
