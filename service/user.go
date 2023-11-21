package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"sync"
	"time"
	"tmp/dao"
	"tmp/helper"
	"tmp/model"
)

var (
	usercountmutex sync.RWMutex
	usercount      int64
)

func UserRegister(c *gin.Context) {
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	pwd := helper.GetMd5(c.PostForm("password"))
	user := &model.User{
		Name:     name,
		Phone:    phone,
		Password: pwd,
		Money:    0,
	}
	var count int64
	dao.DB.Model(model.User{}).Where("phone=?", phone).Count(&count)
	if count != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该手机号已注册",
		})
		return
	}
	err := dao.DB.Create(&user).Error
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

// 用户上传订单
func UserUploadOrder(c *gin.Context) {
	userid := c.PostForm("userid")
	money := c.PostForm("money")
	floatmoney, err := strconv.ParseFloat(money, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "金额格式错误",
		})
		return
	}
	createtime := time.Now()
	startplace := c.PostForm("startplace")
	endplace := c.PostForm("endplace")
	status := "已创建"
	order := &model.Order{
		UserId:     userid,
		Money:      floatmoney,
		CreateTime: createtime,
		StartPlace: startplace,
		EndPlace:   endplace,
		Status:     status,
	}
	bytes := helper.OrderStruct2Bytes(*order)
	err = helper.Publish("order", bytes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发布消息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

// 用户查询订单状态
func UserQueryOrder(c *gin.Context) {
	userid := c.PostForm("userid")
	result, err := dao.Rdb.Get(context.Background(), "user"+userid).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "等待接单中",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询订单状态出错",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  result,
	})
	return
}
