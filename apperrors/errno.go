package apperrors

type RetCode int32

const (
	RetCode_SUCCESS RetCode = 0
	RetCode_Base            = 51001000
	//基本错误码(01~20)
	RetCode_ERR_BEGIN              RetCode = RetCode_Base + 1
	RetCode_ERR_PARAM              RetCode = RetCode_Base + 2
	RetCode_ERR_TIMEOUT            RetCode = RetCode_Base + 3
	RetCode_ERR_RATE_LIMIT_EXCCEED RetCode = RetCode_Base + 4
	RetCode_ERR_VERSION_TOO_OLD    RetCode = RetCode_Base + 5
	RetCode_ERR_SERVER_INNER       RetCode = RetCode_Base + 6
	RetCode_ERR_OVERLOAD           RetCode = RetCode_Base + 7
	RetCode_ERR_FLOW_NOT_MATCH     RetCode = RetCode_Base + 8
	RetCode_ERR_TYPE_ASSERT        RetCode = RetCode_Base + 9
	RetCode_ERR_MARSH              RetCode = RetCode_Base + 10

	//redis相关错误(21~40)
	RetCode_ERR_CACHE_INIT  RetCode = RetCode_Base + 21
	RetCode_ERR_CACHE_GET   RetCode = RetCode_Base + 22
	RetCode_ERR_CACHE_MISS  RetCode = RetCode_Base + 23
	RetCode_ERR_CACHE_SET   RetCode = RetCode_Base + 24
	RetCode_ERR_CACHE_SETE  RetCode = RetCode_Base + 25
	RetCode_ERR_CACHE_TTL   RetCode = RetCode_Base + 26
	RetCode_ERR_CACHE_DEL   RetCode = RetCode_Base + 27
	RetCode_ERR_CACHE_PUB   RetCode = RetCode_Base + 28
	RetCode_ERR_CACHE_MGET  RetCode = RetCode_Base + 29
	RetCode_ERR_CACHE_SETX  RetCode = RetCode_Base + 30
	RetCode_ERR_CACHE_FLUSH RetCode = RetCode_Base + 31
	RetCode_ERR_CACHE_LPUSH RetCode = RetCode_Base + 32
	RetCode_ERR_CACHE_MSET  RetCode = RetCode_Base + 33
	RetCode_ERR_CACHE_PING  RetCode = RetCode_Base + 34

	//DB相关错误(41~60)
	RetCode_ERR_DB_NOT_READY RetCode = RetCode_Base + 41
	RetCode_ERR_DB_SERVER    RetCode = RetCode_Base + 42

	//mq相关错误(81~110)
	RetCode_ERR_MQ_CONSUMER_MSG RetCode = RetCode_Base + 81
	RetCode_ERR_MQ_PUBLISH_MSG  RetCode = RetCode_Base + 82
	RetCode_ERR_MQ_CREATE_QUEUE RetCode = RetCode_Base + 83
	RetCode_ERR_MQ_DELETE_QUEUE RetCode = RetCode_Base + 84
	RetCode_ERR_MQ_RECOSUME     RetCode = RetCode_Base + 85

	RetCode_ERR_END RetCode = RetCode_Base + 999
)

var RetCodeMsg = map[RetCode]string{
	RetCode_SUCCESS: "成功",

	//基本错误码0~20
	RetCode_ERR_BEGIN:              "错误码开始",
	RetCode_ERR_PARAM:              "请求参数异常",
	RetCode_ERR_TIMEOUT:            "请求超时",
	RetCode_ERR_RATE_LIMIT_EXCCEED: "接口访问太过频繁",
	RetCode_ERR_VERSION_TOO_OLD:    "接口版本太低",
	RetCode_ERR_SERVER_INNER:       "服务内部异常",
	RetCode_ERR_OVERLOAD:           "服务过载",
	RetCode_ERR_FLOW_NOT_MATCH:     "连接信息不匹配",
	RetCode_ERR_TYPE_ASSERT:        "类型转换失败",
	RetCode_ERR_MARSH:              "数据转换json失败",

	//redis错误(21~40)
	RetCode_ERR_CACHE_INIT:  "redis初始化失败",
	RetCode_ERR_CACHE_GET:   "redis读取失败",
	RetCode_ERR_CACHE_MISS:  "redis数据不存在",
	RetCode_ERR_CACHE_SET:   "redis保存数据失败",
	RetCode_ERR_CACHE_SETE:  "redis保存数据失败",
	RetCode_ERR_CACHE_TTL:   "redis设置key生命周期失败",
	RetCode_ERR_CACHE_DEL:   "redis删除数据失败",
	RetCode_ERR_CACHE_PUB:   "redis发布失败",
	RetCode_ERR_CACHE_MGET:  "redis批量获取失败",
	RetCode_ERR_CACHE_SETX:  "redis设置过期时间失败",
	RetCode_ERR_CACHE_FLUSH: "redis刷新失败",
	RetCode_ERR_CACHE_LPUSH: "redis列表插入元素失败",
	RetCode_ERR_CACHE_MSET:  "redis批量写入失败",
	RetCode_ERR_CACHE_PING:  "redis Ping失败",

	//数据库错误(41~60)
	RetCode_ERR_DB_NOT_READY: "DB访问未初始化",
	RetCode_ERR_DB_SERVER:    "DB异常",

	//消息队列错误(81~100)
	RetCode_ERR_MQ_CONSUMER_MSG: "MQ消息获取失败",
	RetCode_ERR_MQ_PUBLISH_MSG:  "MQ发布消息失败",
	RetCode_ERR_MQ_CREATE_QUEUE: "创建队列失败",
	RetCode_ERR_MQ_DELETE_QUEUE: "删除队列失败",
	RetCode_ERR_MQ_RECOSUME:     "MQ消息回退失败",

	RetCode_ERR_END: "错误码结束",
}
