package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"tmp/dao"
	"tmp/helper"
	"tmp/model"
)

// 司机是否继续订阅信号
var (
	stopctx  context.Context
	stopfunc context.CancelFunc
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
	time.Sleep(300 * time.Second)
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

// 初始化取消订阅信号
func InitStopCtx() {
	stopctx, stopfunc = context.WithCancel(context.Background())
}
