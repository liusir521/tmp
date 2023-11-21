package service

import (
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"net/http"
	"strconv"
	"sync"
	"tmp/dao"
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

// 更改司机状态为工作
func DriverWork(c *gin.Context) {
	id := c.PostForm("id")
	err := dao.DB.Model(&model.Driver{}).Where("id=?", id).Update("isworking", true).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新状态出错:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新状态成功",
	})
}

// 更改司机状态为休息
func DriverNotWork(c *gin.Context) {
	id := c.PostForm("id")
	err := dao.DB.Model(&model.Driver{}).Where("id=?", id).Update("isworking", false).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新状态出错:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新状态成功",
	})
}

// 开始接单
func DriverStart(c *gin.Context) {
	driveridstr := c.PostForm("driverid")
	carid := c.PostForm("carid")
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("order", "generalorder", config)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订阅消息出错:" + err.Error(),
		})
	}
	driverid, _ := strconv.ParseInt(driveridstr, 10, 64)
	wg := sync.WaitGroup{}
	drivermessage := &messagehandler.DriverMessageHandler{
		DriverId: driverid,
		CarId:    carid,
		Wg:       wg,
	}
	wg.Add(1)
	consumer.AddHandler(drivermessage)
	wg.Wait()
}
