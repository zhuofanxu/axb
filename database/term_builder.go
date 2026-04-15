package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/zhuofanxu/axb/http/dto"
)

// termPair 将 SQL 片段与其参数绑定，避免索引错位
type termPair struct {
	sql      string
	args     []interface{}
	termType string // "and" / "or"，与前一个条件的逻辑关系
}

// BuildTermConditions 构建 terms 查询条件
func BuildTermConditions(db *gorm.DB, terms []dto.TermGroup) *gorm.DB {
	if len(terms) == 0 {
		return db
	}

	for _, group := range terms {
		if len(group.Terms) == 0 {
			continue
		}

		var pairs []termPair
		for _, term := range group.Terms {
			sql, args := buildTermCondition(term)
			if sql == "" {
				continue
			}
			pairs = append(pairs, termPair{sql: sql, args: args, termType: strings.ToLower(term.Type)})
		}

		if len(pairs) == 0 {
			continue
		}

		if len(pairs) == 1 {
			db = db.Where(pairs[0].sql, pairs[0].args...)
			continue
		}

		// 按 OR/AND 分组：type=="or" 的条件与上一个条件合并为一个 OR 段，
		// type=="and" 时开启新段，各段之间用 AND 连接。
		// 例：[A(and), B(or), C(and), D(or), E(and)] => (A OR B) AND (C OR D) AND E
		type segment struct {
			sql  string
			args []interface{}
		}

		var segments []segment
		pendingSQL := pairs[0].sql
		pendingArgs := append([]interface{}{}, pairs[0].args...)

		for _, p := range pairs[1:] {
			if p.termType == "or" {
				pendingSQL += " OR " + p.sql
				pendingArgs = append(pendingArgs, p.args...)
			} else {
				segments = append(segments, segment{sql: pendingSQL, args: pendingArgs})
				pendingSQL = p.sql
				pendingArgs = append([]interface{}{}, p.args...)
			}
		}
		segments = append(segments, segment{sql: pendingSQL, args: pendingArgs})

		var sqlParts []string
		var allArgs []interface{}
		for _, seg := range segments {
			if strings.Contains(seg.sql, " OR ") {
				sqlParts = append(sqlParts, "("+seg.sql+")")
			} else {
				sqlParts = append(sqlParts, seg.sql)
			}
			allArgs = append(allArgs, seg.args...)
		}

		db = db.Where(strings.Join(sqlParts, " AND "), allArgs...)
	}

	return db
}

// 构建单个查询条件
func buildTermCondition(term dto.Term) (string, []interface{}) {
	if term.Column == "" || term.TermType == "" {
		return "", nil
	}

	// 转换驼峰为蛇形命名
	column := camelToSnake(term.Column)
	termType := strings.ToLower(term.TermType)

	switch termType {
	case "like":
		// LIKE 查询
		if value, ok := term.Value.(string); ok {
			// 如果值中没有 %，自动添加
			if !strings.Contains(value, "%") {
				value = "%" + value + "%"
			}
			return fmt.Sprintf("%s LIKE ?", column), []interface{}{value}
		}

	case "eq", "=":
		// 等于
		return fmt.Sprintf("%s = ?", column), []interface{}{term.Value}

	case "neq", "!=":
		// 不等于
		return fmt.Sprintf("%s != ?", column), []interface{}{term.Value}

	case "in":
		// IN 查询
		return fmt.Sprintf("%s IN (?)", column), []interface{}{term.Value}

	case "nin", "not_in":
		// NOT IN 查询
		return fmt.Sprintf("%s NOT IN (?)", column), []interface{}{term.Value}

	case "gt", ">":
		// 大于
		return fmt.Sprintf("%s > ?", column), []interface{}{term.Value}

	case "gte", ">=":
		// 大于等于
		return fmt.Sprintf("%s >= ?", column), []interface{}{term.Value}

	case "lt", "<":
		// 小于
		return fmt.Sprintf("%s < ?", column), []interface{}{term.Value}

	case "lte", "<=":
		// 小于等于
		return fmt.Sprintf("%s <= ?", column), []interface{}{term.Value}

	case "is_null":
		// IS NULL
		return fmt.Sprintf("%s IS NULL", column), []interface{}{}

	case "is_not_null":
		// IS NOT NULL
		return fmt.Sprintf("%s IS NOT NULL", column), []interface{}{}
	}

	return "", nil
}

// 构建排序条件
func BuildSortConditions(db *gorm.DB, sorts []dto.SortItem) *gorm.DB {
	if len(sorts) == 0 {
		return db
	}

	for _, sort := range sorts {
		if sort.Name == "" {
			continue
		}

		// 转换驼峰为蛇形命名
		column := camelToSnake(sort.Name)
		order := strings.ToUpper(sort.Order)

		// 验证排序方向
		if order != "ASC" && order != "DESC" {
			order = "ASC" // 默认升序
		}

		db = db.Order(fmt.Sprintf("%s %s", column, order))
	}

	return db
}
