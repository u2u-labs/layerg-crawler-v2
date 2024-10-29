package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

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

type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
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

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func QueryWithDynamicFilter[T any](db *sql.DB, tableName string, limit int, offset int, filterConditions map[string]string) ([]T, error) {
	// Set up Squirrel with PostgreSQL placeholder format
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Create a Squirrel query builder for the items table
	queryBuilder := psql.Select("*").From(tableName)

	// Apply dynamic filters
	for column, value := range filterConditions {
		queryBuilder = queryBuilder.Where(squirrel.Eq{column: value})
	}

	// Apply pagination
	queryBuilder = queryBuilder.Limit(uint64(limit)).Offset(uint64(offset))

	// Convert the query to SQL
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL query: %w", err)
	}

	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Get the column names from the result
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %w", err)
	}

	// Create a slice to store the results
	var items []T

	// Create a new item of type T to get its structure
	var item T
	itemType := reflect.TypeOf(item)
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	// Create a map of column names to struct field indices
	columnMap := make(map[string]int)
	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)

		// First check for db tag
		tag := field.Tag.Get("db")
		if tag != "" {
			columnMap[tag] = i
			continue
		}

		// If no db tag, convert field name to snake_case
		snakeCaseName := toSnakeCase(field.Name)
		columnMap[snakeCaseName] = i
	}

	// Scan rows into the slice of items
	for rows.Next() {
		// Create a new item for each row
		itemValue := reflect.New(itemType).Elem()

		// Create a slice of interface{} to hold the row values
		scanArgs := make([]interface{}, len(columns))
		for i, colName := range columns {
			if fieldIndex, ok := columnMap[colName]; ok {
				scanArgs[i] = itemValue.Field(fieldIndex).Addr().Interface()
			} else {
				// Handle columns that don't map to struct fields
				var placeholder interface{}
				scanArgs[i] = &placeholder
			}
		}

		// Scan the row into the struct fields
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Append the item to our slice
		items = append(items, itemValue.Interface().(T))
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return items, nil
}

// CountItems counts the number of items in the database based on dynamic filters.
func CountItemsWithFilter(db *sql.DB, tableName string, filterConditions map[string]string) (int, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Create a Squirrel query builder for counting items
	queryBuilder := psql.Select("COUNT(*)").From(tableName)

	// Apply dynamic filters
	for column, value := range filterConditions {
		fmt.Println(column, value)
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
