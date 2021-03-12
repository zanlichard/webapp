package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
)

//communication data format by rabbitmq
type MQMsgFormat interface {
	ObjToMsg() string
	MsgToObj(dataByte []byte) error
}

//rabbitmq operation set
type MQOperation interface {
	//获取可用的channel
	RefreshConnectionAndChannel() (channel *amqp.Channel, err error)

	//创建一个queue队列
	CreateQueue(queue string) error

	//删除一个queue队列
	DeleteQueue(queue string) error

	//发布消息到队列
	PublishQueue(exchange string, routeKey string, body string) error

	//取出消息消费
	ConsumeQueue(context.Context, Receiver) error

	//退回消息
	ReConsume(exchange string, queue string, msg string) error

	//统计正在队列中准备且还未消费的数据
	GetReadyCount(queue string) (int, error)

	//获取到队列中正在消费的数据，这里指的是正在有多少数据被消费
	GetConsumCount(queue string) (int, error)

	//清理队列
	ClearQueue(queue string) (string, error)
}

//消费者定义
type Receiver interface {
	OnError(error)         // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
	OnReceive([]byte) bool // 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
	QueueName() string     // 获取接收者需要监听的队列
	RouterKey() string     // 这个队列绑定的路由
}
