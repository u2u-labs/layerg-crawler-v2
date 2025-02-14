package core

import (
	"fmt"
	"strings"
)

func BuildQuery(entity string, args map[string]interface{}) string {
	base := fmt.Sprintf("SELECT * FROM %s", strings.ToLower(entity))
	var conds []string
	if where, ok := args["where"].(map[string]interface{}); ok {
		for k, v := range where {
			conds = append(conds, fmt.Sprintf("%s = '%v'", k, v))
		}
	}
	if len(conds) > 0 {
		base += " WHERE " + strings.Join(conds, " AND ")
	}
	if order, ok := args["order"].(string); ok {
		base += " ORDER BY " + order
	}
	if limit, ok := args["limit"].(int); ok {
		base += fmt.Sprintf(" LIMIT %d", limit)
	}
	if page, ok := args["page"].(int); ok {
		if limit, ok := args["limit"].(int); ok {
			base += fmt.Sprintf(" OFFSET %d", (page-1)*limit)
		}
	}
	return base
}
