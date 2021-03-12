package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"webapp/logger"
)

type MQData struct {
	From        string `json:"from"`         //转出地址
	To          string `json:"to"`           //转入地址
	Amount      string `json:"amount"`       //金额
	Fee         string `json:"fee"`          //手续费 eth
	BusinessId  string `json:"business_id"`  //业务id 用于审核前区分
	TxId        string `json:"tx_id"`        //钱包转账生成的txid 全局唯一 用于后续审核通知
	Uid         int32  `json:"uid"`          //用户userid
	Status      int    `json:"status"`       //1：审核中 2:成功 3:失败
	ToUid       int32  `json:"to_uid"`       //投票专用
	AddressType int32  `json:"address_type"` //地址类型 1钱包地址，2投票地址，3节点地址
	CoinType    string `json:"coin_type"`    //1:speed 2:eth 3:fn 4:usdt
}

//实现发送者,将对象转换为json字符串
func (m *MQData) ObjToMsg() string {
	mJson, err := json.Marshal(m)
	if err != nil {
		logger.Logger.Error("json.Marshal(m) %+v", err)
	}
	return string(mJson)
}

// 实现接收者,将对象数据格式化为json
func (m *MQData) MsgToObj(dataByte []byte) error {
	return json.Unmarshal(dataByte, m)
}

type TestReceiver struct {
	Queue  string
	Router string // 这个队列绑定的路由
}

//实现消费者的接口
func (receiver *TestReceiver) OnError(err error) {
	fmt.Printf("message consumer failed for:%+v", err)
}

func (receiver *TestReceiver) OnReceive(msg []byte) bool {
	rd := new(MQData)
	rd.MsgToObj(msg)
	fmt.Printf("consume data:%+v", rd)
	return true
}

func (receiver *TestReceiver) RouterKey() string {
	return receiver.Router
}

func (receiver *TestReceiver) QueueName() string {
	return receiver.Queue
}

func TestRabbitMqWrite(t *testing.T) {
	mq, err := NewRabbitMq("devops", "devops", "192.168.163.129", 5672, "/order_host")
	if err != nil {
		t.Errorf("init rabbitmq failed for:%+v", err)
		return
	}
	data := &MQData{
		From:        "lichard",
		To:          "your",
		Amount:      "1000",
		Fee:         "1.0",
		BusinessId:  "1000001",
		TxId:        "0x78gdgdsbde",
		Uid:         1387995,
		Status:      0,
		ToUid:       781479142,
		AddressType: 3,
		CoinType:    "usdt",
	}

	err = mq.PublishQueue("", "test", data.ObjToMsg())
	if err != nil {
		t.Errorf("publish msg failed for:%+v", err)
		return
	}

}

func TestRabbitMqRead(t *testing.T) {
	mq, err := NewRabbitMq("devops", "devops", "192.168.163.129", 5672, "/order_host")
	if err != nil {
		t.Errorf("init rabbitmq failed for:%+v", err)
		return
	}
	consumer := TestReceiver{
		Queue:  "test",
		Router: "",
	}
	err = mq.ConsumeQueue(context.Background(), &consumer)
	if err != nil {
		t.Errorf("consume queue failed for:%+v", err)
		return
	}

}

func TestRabbitMqCreateQueue(t *testing.T) {
	mq, err := NewRabbitMq("devops", "devops", "192.168.163.129", 5672, "/order_host")
	if err != nil {
		t.Errorf("init rabbitmq failed for:%+v", err)
		return
	}

	err = mq.CreateQueue("test")
	if err != nil {
		t.Errorf("create queue failed for:%+v", err)
		return
	}
}

func TestRabbitMqDelQueue(t *testing.T) {
	mq, err := NewRabbitMq("devops", "devops", "192.168.163.129", 5672, "/order_host")
	if err != nil {
		t.Errorf("init rabbitmq failed for:%+v", err)
		return
	}

	err = mq.DeleteQueue("test")
	if err != nil {
		t.Errorf("consume queue delete failed for:%+v", err)
		return
	}
}
