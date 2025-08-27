package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateInsertQuery(table string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var column, placeholder string

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			if column != "" {
				column += ", "
				placeholder += ", "
			}
			column += dbTag
			placeholder += "?"
		}

	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)",table, column, placeholder)
}

func GetStructValues(model interface{}) []interface{} {
	modelVal := reflect.ValueOf(model)
	modelType := modelVal.Type()

	values := []interface{}{}

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelVal.Field(i).Interface())
		}
	}
	return values
}
