package messagehandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"net/http"
	"time"
	"tmp/dao"
	"tmp/global"
	"tmp/model"
)

type DriverMessageHandler struct {
	DriverId  int64
	CarId     string
	Consummer *nsq.Consumer
	Res       *gin.Context
}

func (driver DriverMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	var order model.Order
	err := json.Unmarshal(m.Body, &order)
	if err != nil {
		return err
	}
	order.Carid = driver.CarId
	order.DriverId = driver.DriverId
	fmt.Println(driver.DriverId, "接到订单", driver.CarId, order)
	order.Status = "已接单"
	orderkey := "order" + string(order.ID)

	// redis分布式锁避免重复消费
	_, err = dao.Rdb.SetNX(context.Background(), orderkey, order, 5*time.Second).Result()
	if err != nil {
		return err
	} else {
		// 将接到的订单信息放入到订单对象池中
		err = global.GlobalOrder.Add(orderkey, order, 0)
		if err != nil {
			return err
		}
		driver.Res.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  order,
		})
		m.Finish()
		driver.Consummer.Stop()
		return nil
	}
}
