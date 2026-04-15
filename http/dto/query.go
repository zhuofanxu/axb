package dto

// QueryRequest 通用查询请求
type QueryRequest struct {
	PaginationParam
	Sorts []SortItem  `json:"sorts"` // 排序条件
	Terms []TermGroup `json:"terms"` // 查询条件组
}

// SortItem 排序项
type SortItem struct {
	Name  string `json:"name"`  // 字段名
	Order string `json:"order"` // asc/desc
}

// TermGroup 条件组（支持嵌套）
type TermGroup struct {
	Terms []Term `json:"terms"` // 条件列表
	Type  string `json:"type"`  // and/or，组内条件的逻辑关系
}

// Term 查询条件
type Term struct {
	Column   string      `json:"column"`   // 字段名
	TermType string      `json:"termType"` // 查询类型：like, eq, neq, in, nin, gt, gte, lt, lte
	Value    interface{} `json:"value"`    // 查询值
	Type     string      `json:"type"`     // and/or，与前一个条件的逻辑关系
	Terms    []Term      `json:"terms"`    // 嵌套条件
}
