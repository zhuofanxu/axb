package database

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

const (
	defaultPageIndex = 0
	defaultPageSize  = 20
)

// BuildQueryCondition builds a GORM query condition based on the provided struct parameters.
func BuildQueryCondition(db *gorm.DB, params interface{}) *gorm.DB {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return db
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("query_field")
		if value.IsZero() || tag == "" {
			continue
		}
		db = db.Where(fmt.Sprintf("%s = ?", tag), value.Interface())
	}
	return db
}

func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func BuildUpdateMap(params interface{}) map[string]interface{} {
	updateMap := make(map[string]interface{})
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return updateMap
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("json")
		// 提取 JSON 标签名（去掉 omitempty 等）
		tag = strings.Split(tag, ",")[0]

		if value.IsZero() || tag == "" || tag == "-" {
			continue
		}
		dbFieldName := camelToSnake(tag)
		updateMap[dbFieldName] = value.Interface()
	}
	return updateMap
}

func PaginateQuery[T any](db *gorm.DB, records *[]T, page, pageSize int) (int64, error) {
	return executePagination(db, db, records, page, pageSize)
}

// PaginateComplexQuery 处理复杂查询分页
//
//goland:noinspection GoUnusedExportedFunction
func PaginateComplexQuery[T any](dataQuery, countQuery *gorm.DB, records *[]T, page, pageSize int) (int64, error) {
	return executePagination(dataQuery, countQuery, records, page, pageSize)
}

func executePagination[T any](dataQuery, countQuery *gorm.DB, records *[]T, page, pageSize int) (int64, error) {
	page, pageSize = normalizePagination(page, pageSize)

	total, err := queryTotal(countQuery)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		*records = []T{} // assign an empty slice if no records found
		return 0, nil
	}
	if err = queryPage(dataQuery, records, page, pageSize); err != nil {
		return 0, WrapDBErr(err)
	}

	return total, nil
}

func queryTotal(db *gorm.DB) (int64, error) {
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, WrapDBErr(err)
	}
	return total, nil
}

func queryPage[T any](db *gorm.DB, records *[]T, page, pageSize int) error {
	offset := page * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(records).Error; err != nil {
		return err
	}
	return nil
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 0 {
		page = defaultPageIndex
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	return page, pageSize
}
