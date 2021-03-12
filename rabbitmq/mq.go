package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"
	. "webapp/logger"

	"github.com/streadway/amqp"
)

const (
	MQ_STATUS_CHECKING = 1
	MQ_STATUS_OK       = 2
	MQ_STATUS_FAIL     = 3
)

type Rabbitmq struct {
	Conn               *amqp.Connection
	Lock               *sync.RWMutex
	err                error
	RabbitmqUsername   string
	RabbitmqPassword   string
	RabbitmqServeraddr string
	RabbitmqServerPort int
	RabbitmqVhost      string
}

//开始创建一个新的rabitmq对象
func NewRabbitMq(Username string, Password string, Serveraddr string, ServerPort int, Vhost string) (*Rabbitmq, error) {
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", Username, Password, Serveraddr, ServerPort, Vhost)
	conn, err := amqp.Dial(RabbitUrl) //默认10s心跳,编码(us-en)
	if err != nil {
		return nil, err
	}
	rabbitmq := new(Rabbitmq)
	rabbitmq.Conn = conn
	rabbitmq.RabbitmqUsername = Username
	rabbitmq.RabbitmqPassword = Password
	rabbitmq.RabbitmqServeraddr = Serveraddr
	rabbitmq.RabbitmqServerPort = ServerPort
	rabbitmq.RabbitmqVhost = Vhost
	rabbitmq.Lock = new(sync.RWMutex)
	return rabbitmq, nil
}

func (rabbitmq *Rabbitmq) refreshConnectionAndChannel() (channel *amqp.Channel, err error) {
	rabbitmq.Lock.Lock()
	defer rabbitmq.Lock.Unlock()
	channel, err = rabbitmq.Conn.Channel()
	if err != nil {
		for {
			rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitmq.RabbitmqUsername, rabbitmq.RabbitmqPassword, rabbitmq.RabbitmqServeraddr, rabbitmq.RabbitmqServerPort,
				rabbitmq.RabbitmqVhost)
			rabbitmq.Conn, err = amqp.Dial(rabbitUrl)
			if err != nil {
				Logger.Info("connect mq error for:%v,retry...", err)
				time.Sleep(5 * time.Second)
			} else {
				return rabbitmq.Conn.Channel()
			}
		}
	}
	return
}

func (rabbitmq *Rabbitmq) CreateQueue(queue string) error {
	ch, err := rabbitmq.refreshConnectionAndChannel()
	defer ch.Close()
	if err != nil {
		Logger.Error("rabbitmq.CreateQueue queue:%s get channel failed for: %+v", queue, err)
		return err
	}
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		Logger.Error("rabbitmq.CreateQueue queue:%s failed for: %+v", queue, err)
		return err
	}
	return nil
}

func (rabbitmq *Rabbitmq) DeleteQueue(queue string) error {
	ch, err := rabbitmq.refreshConnectionAndChannel()
	defer ch.Close()
	if err != nil {
		Logger.Error("rabbitmq.DeleteQueue queue:%s get channel failed for: %+v", queue, err)
		return err
	}
	_, err = ch.QueueDelete(
		queue, // name
		false, // IfUnused
		false, // ifEmpty
		true,  // noWait
	)
	if err != nil {
		Logger.Error("rabbitmq.DeleteQueue queue:%s failed for: %+v", queue, err)
		return err
	}
	return nil
}

func (rabbitmq *Rabbitmq) PublishQueue(exchange string, routeKey string, body string) error {
	ch, err := rabbitmq.refreshConnectionAndChannel()
	defer ch.Close()
	if err != nil {
		Logger.Error("rabbitmq.PublishQueue exchange:%s key:%s get channel failed for: %+v", exchange, routeKey, err)
		return err
	}
	err = ch.Publish(
		exchange, // exchange
		routeKey, // routing key
		false,    // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		Logger.Error("rabbitmq.PublishQueue exchange:%s key:%s publish mesage failed for: %+v", exchange, routeKey, err)
		return err
	}
	return nil
}

func (rabbitmq *Rabbitmq) ConsumeQueue(ctx context.Context, receiver Receiver) error {
	ch, err := rabbitmq.refreshConnectionAndChannel()
	if err != nil {
		Logger.Error("rabbitmq.ConsumeQueue get channel failed for:%+v", err)
		receiver.OnError(err)
		return err
	}
	defer ch.Close()
	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		Logger.Error("rabbitmq.ConsumeQueue prefetch failed for:%+v", err)
		receiver.OnError(err)
		return err
	}
	queue := receiver.QueueName()
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		Logger.Error("rabbitmq.ConsumeQueue queue:%s consume failed for:%+v", queue, err)
		return err
	}

	go func() {
		for d := range msgs {
			receiver.OnReceive(d.Body)
			if err != nil {
				Logger.Error("rabbitmq.ConsumerData queue:%s failed for:%+v", queue, err)
				receiver.OnError(err)
			}
			d.Ack(false)
			//d.Nack(false, true)
		}
	}()
	return nil
}

func (rabbitmq *Rabbitmq) ReConsume(exchange string, queue string, msg string) error {
	merr := rabbitmq.PublishQueue(exchange, queue, msg)
	if merr != nil {
		Logger.Error("rabbitmq.ReConsume PublishQueue exchange:%s queue:%s failed for:%+v", exchange, queue, merr)
		return merr
	}
	return nil
}

func (rabbitmq *Rabbitmq) GetReadyCount(queue string) (int, error) {
	count := 0
	ch, err := rabbitmq.refreshConnectionAndChannel()
	defer ch.Close()
	if err != nil {
		Logger.Error("rabbitmq.GetReadyCount queue:%s get channel failed for: %+v", queue, err)
		return count, err
	}
	state, err := ch.QueueInspect(queue)
	if err != nil {
		Logger.Error("rabbitmq.GetReadyCount queue:%s Inspect failed for: %+v", queue, err)
		return count, err
	}
	return state.Messages, nil
}

func (rabbitmq *Rabbitmq) GetConsumCount(queue string) (int, error) {
	count := 0
	ch, err := rabbitmq.refreshConnectionAndChannel()
	defer ch.Close()
	if err != nil {
		Logger.Error("rabbitmq.GetConsumCount queue:%s get channel failed for: %+v", queue, err)
		return count, err
	}
	state, err := ch.QueueInspect(queue)
	if err != nil {
		Logger.Error("rabbitmq.GetConsumCount queue:%s Inspect failed for: %+v", queue, err)
		return count, err
	}
	return state.Consumers, nil
}

func (rabbitmq *Rabbitmq) ClearQueue(queue string) (string, error) {
	ch, err := rabbitmq.refreshConnectionAndChannel()
	defer ch.Close()
	if err != nil {
		Logger.Error("rabbitmq.ClearQueue queue:%s get channel failed for: %+v", queue, err)
		return "", err
	}
	_, err = ch.QueuePurge(queue, false)
	if err != nil {
		Logger.Error("rabbitmq.ClearQueue queue:%s purge failed for: %+v", queue, err)
		return "", err
	}
	return "Delete queue success", nil
}
