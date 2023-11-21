package messagehandler

import (
	"context"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
	"tmp/dao"
	"tmp/global"
	"tmp/model"
)

// 通过
type DriverMessageHandler struct {
	DriverId int64
	CarId    string
	Wg       sync.WaitGroup
}

func (driver DriverMessageHandler) HandleMessage(m *nsq.Message) error {
	defer driver.Wg.Done()
	if len(m.Body) == 0 {
		return nil
	}
	var order model.Order
	err := json.Unmarshal(m.Body, &order)
	if err != nil {
		return err
	}
	order.Carid = driver.CarId
	order.DriverId = string(driver.DriverId)
	order.Status = "已接单"
	orderkey := "order" + string(order.ID)

	// redis分布式锁避免重复消费
	err = dao.Rdb.SetNX(context.Background(), orderkey, order, 5*time.Second).Err()
	//err = dao.Rdb.Set(context.Background(), "user"+order.UserId, order, 0).Err()
	if err != nil {
		return err
	}

	// 将接到的订单信息放入到订单对象池中
	err = global.GlobalOrder.Add(orderkey, order, 0)
	if err != nil {
		return err
	}
	return nil
}
