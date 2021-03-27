package code

const (
	//SUCCESS              = 200
	SUCCESS                        = 0 //   开发完了daving让改为0
	INVALID_PARAMS                 = 400
	ID_NOT_EMPTY                   = 4001
	ERROR_TOKEN_EMPTY              = 4002
	ERROR_TOKEN_INVALID            = 4003
	ERROR_TOKEN_EXPIRE             = 4004
	ERROR_USER_NOT_EXIST           = 4005
	ERROR                          = 500
	ERROR_DATA_NOT_EXIST           = 5001
	ERROR_CONFIG_PARSE             = 5003
	ERROR_SIGN_FIELD_NO_EXIST      = 5004
	ERROR_SIGN                     = 5005
	ERRO_SERVICE_ID_FIELD_NO_EXIST = 5006
	ERROR_DENY_SERVICE_ID          = 5007
	ERROR_LOST_SIGN_DATA           = 5008
)

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	INVALID_PARAMS:                 "请求参数错误",
	ID_NOT_EMPTY:                   "ID为空",
	ERROR_TOKEN_EMPTY:              "token为空",
	ERROR_TOKEN_INVALID:            "token无效",
	ERROR_TOKEN_EXPIRE:             "token过期",
	ERROR_USER_NOT_EXIST:           "用户不存在",
	ERROR:                          "服务内部错误",
	ERROR_DATA_NOT_EXIST:           "记录不存在",
	ERROR_CONFIG_PARSE:             "解析配置出错",
	ERROR_SIGN_FIELD_NO_EXIST:      "没有签名字段",
	ERROR_SIGN:                     "签名错误",
	ERRO_SERVICE_ID_FIELD_NO_EXIST: "没有服务ID",
	ERROR_DENY_SERVICE_ID:          "服务未授权",
	ERROR_LOST_SIGN_DATA:           "没有签名数据",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
