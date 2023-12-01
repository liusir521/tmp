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

type UserMessageHandler struct {
	UserId    int64
	CarId     string
	Consummer *nsq.Consumer
	Res       *gin.Context
}

func (user UserMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	var order model.Order
	err := json.Unmarshal(m.Body, &order)
	if err != nil {
		return err
	}
	fmt.Println(user.UserId, ":顺风车:", order)
	orderkey := "order" + string(order.DriverId) + string(user.UserId)

	// redis分布式锁避免重复消费
	_, err = dao.Rdb.SetNX(context.Background(), orderkey, order, 5*time.Second).Result()
	if err != nil {
		return err
	} else {
		_, b := global.GlobalTogetherOrder.Get("driver" + string(order.DriverId))
		// 找到说明司机未发车，可以消费
		if b {
			// 将乘客添加到顺风车订单中
			AddUser("driver"+string(order.DriverId), user.UserId)
			fmt.Println(global.TogetherUsers.Load("driver" + string(order.DriverId)))
			user.Res.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  order,
			})
			m.Finish()
			user.Consummer.Stop()
			return nil
		} else {
			// 未找到说明司机已经发车，将此订单消费并消费下一个订单
			m.Finish()
			return nil
		}
	}
}

// 将用户添加到顺风车乘客中
func AddUser(key string, val int64) {
	value, _ := global.TogetherUsers.Load(key)
	userids := value.([]int64)
	swap := global.TogetherUsers.CompareAndSwap(key, userids, append(userids, val))
	if !swap {
		AddUser(key, val)
	}
}
