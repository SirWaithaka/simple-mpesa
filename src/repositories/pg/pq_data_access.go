package pg

import (
	"fmt"
	"strings"
)

// GetBy builds a select statement used to get records from db using filter
// SELECT * FROM <tablename> WHERE filter1=filtervalue
func GetBy(tableName string, filter map[string]interface{}) (queryStr string, dataFields []interface{}){
	queryStr = fmt.Sprintf("SELECT * FROM %s WHERE", tableName)
	for k, v := range filter {
		queryStr += fmt.Sprintf(" %s=?,", k)
		dataFields = append(dataFields, v)
	}
	queryStr = strings.TrimSuffix(queryStr, ",")

	return queryStr, dataFields
}

// GetLimitBy builds a select statement used to get records from db using limit and filter
func GetLimitBy(tableName string, limit uint, filter map[string]interface{}) (queryStr string, dataFields []interface{}) {
	queryStr = fmt.Sprintf("SELECT * FROM %s WHERE", tableName)
	for k, v := range filter {
		queryStr += fmt.Sprintf(" %s=?,", k)
		dataFields = append(dataFields, v)
	}
	queryStr = strings.TrimSuffix(queryStr, ",")

	// add limit
	queryStr += fmt.Sprintf(" LIMIT %s", limit)

	return queryStr, dataFields
}
