package service

import (
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"net/http"
	"strconv"
	"tmp/conf"
	"tmp/dao"
	"tmp/global"
	"tmp/helper"
	"tmp/messagehandler"
	"tmp/model"
)

// 创建司机
func DriverRegister(c *gin.Context) {
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	carid := c.PostForm("carid")
	pwd := helper.GetMd5(c.PostForm("password"))
	driver := &model.Driver{
		Name:      name,
		Phone:     phone,
		Password:  pwd,
		Money:     0,
		IsWorking: false,
		RunCount:  0,
		CarId:     carid,
	}
	var count int64
	dao.DB.Model(model.Driver{}).Where("phone=?", phone).Count(&count)
	// 检查手机号是否已创建
	if count != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该手机号已注册",
		})
		return
	}
	err := dao.DB.Create(&driver).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create User Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
	return
}

// 将司机状态改为工作
// 添加司机到缓存对象池，且状态为工作
func DriverWork(c *gin.Context) {
	driverid := c.PostForm("driverid")
	atoi, _ := strconv.Atoi(driverid)
	err := dao.DB.Model(&model.Driver{}).Where("id=?", atoi).Update("isworking", true).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新状态出错:" + err.Error(),
		})
		return
	}
	var curdriver model.Driver
	err = dao.DB.Where("id=?", driverid).Take(&curdriver).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "司机信息出错:" + err.Error(),
		})
		return
	}

	// 将司机对象添加到缓存池中
	err = global.GlobalDriver.Add("driver"+driverid, curdriver, 0)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "司机添加出错:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新状态成功",
	})
}

// 更改司机状态为休息
// 并从缓存池中删除相关对象
func DriverNotWork(c *gin.Context) {
	driverid := c.PostForm("driverid")
	atoi, _ := strconv.Atoi(driverid)
	err := dao.DB.Model(&model.Driver{}).Where("id=?", atoi).Update("isworking", false).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新状态出错:" + err.Error(),
		})
		return
	}
	// 删除缓存池中的司机对象
	global.GlobalDriver.Delete("driver" + driverid)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新状态成功",
	})
}

// 开始接单
func DriverStart(c *gin.Context) {
	driveridstr := c.PostForm("driverid")
	driverid, _ := strconv.ParseInt(driveridstr, 10, 64)
	carid := c.PostForm("carid")
	config := nsq.NewConfig()
	// 创建消费者
	consumer, err := nsq.NewConsumer("order", "generalorder", config)
	drivermessage := &messagehandler.DriverMessageHandler{
		DriverId:  driverid,
		CarId:     carid,
		Consummer: consumer,
		Res:       c,
	}
	consumer.AddHandler(drivermessage)
	defer consumer.Stop()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订阅消息出错:" + err.Error(),
		})
	}
	// 连接消息队列
	err = consumer.ConnectToNSQD(conf.NsqdAddr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "连接消息队列出错:" + err.Error(),
		})
	}
	<-consumer.StopChan
}

// 完成订单
func FinshOrder(c *gin.Context) {
	userid := "user" + c.PostForm("userid")
	orderid := "order" + c.PostForm("orderid")
	// 将用户从其缓存池中删除
	global.GlobalUser.Delete(userid)
	get, b := global.GlobalOrder.Get(orderid)
	if b {
		order := get.(model.Order)
		order.Status = "已完成"
		// 将订单信息更新到数据库
		err := dao.DB.Save(&order).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "订单更新失败" + err.Error(),
			})
		}
		// 将订单从订单缓存池中删除
		global.GlobalOrder.Delete(orderid)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订单信息出错",
		})
	}
}
