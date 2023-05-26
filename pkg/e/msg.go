package e

var MsgFlags = map[uint]string{
	Success:      "操作成功",
	Error:        "操作失败",
	InvalidParam: "请求参数错误",

	ErrorAuthCheckTokenFail:    "token 错误",
	ErrorAuthCheckTokenTimeout: "token 过期",
}

// GetMsg 获取状态码对应信息
func GetMsg(code uint) string {
	if msg, ok := MsgFlags[code]; ok {
		return msg
	}
	return MsgFlags[Error]

}
