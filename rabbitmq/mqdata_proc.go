package rabbitmq
import(
   "github.com/streadway/amqp"
	"context"
)




type MQMsgFormat interface {
	ObjToMsg() string
	MsgToObj(dataByte []byte) error
}


type MQOperation interface {
	//获取可用的channel
	RefreshConnectionAndChannel() (channel *amqp.Channel,err error)

	//创建一个queue队列
	CreateQueue(queue string) error

	//删除一个queue队列
	DeleteQueue(queue string) error

	//发布消息到队列
	PublishQueue(exchange string,routeKey string, body string) error

	//取出消息消费
	ConsumeQueue(ctx context.Context,queue string,ElementProc func (data *MQData) (error) ) error

	//退回消息
	ReConsume(exchange string,queue string) error

	//统计正在队列中准备且还未消费的数据
	GetReadyCount(queue string) (int, error)

	//获取到队列中正在消费的数据，这里指的是正在有多少数据被消费
	GetConsumCount(queue string) (int, error)

	//清理队列
	ClearQueue(queue string) (string, error)

}
