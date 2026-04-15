package errx

// custom business error codes
const (
	CodeSuccess             = 0      // 成功
	CodeSystemError         = 100000 // 系统内部错误
	CodeParamError          = 100001 // 参数错误
	CodeUnauthorized        = 100002 // 未授权
	CodeForbidden           = 100003 // 禁止访问
	CodeNotFound            = 100004 // 资源不存在
	CodeAlreadyExists       = 100005 // 资源已存在
	CodeCanceled            = 100006 // 操作被取消
	CodeTimout              = 100007 // 操作超时
	CodeRelatedRecordsExist = 100008 // 存在关联记录，无法删除

	CodeInvalidToken  = 200001 // 无效的Token
	CodeTokenExpired  = 200002 // Token过期
	CodeUserNotFound  = 200003 // 用户不存在
	CodePasswordError = 200004 // 密码错误

	CodeNoChanged = 300000 // 无数据变更
)

// 错误码对应的消息
var codeMessages = map[int]string{
	CodeSuccess:             "成功",
	CodeSystemError:         "系统内部错误",
	CodeParamError:          "参数错误",
	CodeUnauthorized:        "未授权",
	CodeForbidden:           "禁止访问",
	CodeNotFound:            "资源不存在",
	CodeAlreadyExists:       "资源已存在",
	CodeNoChanged:           "无数据变更",
	CodeCanceled:            "操作被取消",
	CodeTimout:              "操作超时",
	CodeRelatedRecordsExist: "存在关联记录，无法删除",

	CodeInvalidToken:  "无效的Token",
	CodeTokenExpired:  "Token已过期",
	CodeUserNotFound:  "用户不存在",
	CodePasswordError: "密码错误",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
