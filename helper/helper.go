package helper

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"tmp/model"
)

func GetMd5(pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
}

// 将字节流转换为订单结构体, 用于消息订阅
func OrderBytes2Struct(data []byte) (order model.Order) {
	var order2 model.Order
	err := json.Unmarshal(data, &order2)
	if err != nil {
		panic(err)
	}
	return order2
}

// 将订单结构体转换为字节流，用于消息发布
func OrderStruct2Bytes(order model.Order) []byte {
	marshal, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}
	return marshal
}

func Publish(topic string, data []byte) error {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("192.168.0.102:4150", config)
	defer producer.Stop()
	if err != nil {
		return err
	}
	err = producer.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}
