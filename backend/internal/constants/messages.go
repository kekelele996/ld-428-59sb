package constants

const (
	MsgOK            = "ok"
	MsgInvalidParams = "请求参数错误"
	MsgUnauthorized  = "未登录或登录已过期"
	MsgForbidden     = "没有操作权限"
	MsgNotFound      = "资源不存在"
	MsgInternalError = "服务器内部错误"
	MsgLoginSuccess  = "登录成功"
	MsgRateLimited   = "请求过于频繁，请稍后再试"

	MsgReservationFull   = "该场次预约人数已满"
	MsgReservationClosed = "展览未开放预约"
	MsgPhoneInvalid      = "手机号格式不正确"
	MsgVisitDateInvalid  = "参观日期不在展期内"
	MsgSessionInvalid    = "参观场次不存在"
)
