package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
	"tmp/dao"
	"tmp/global"
	"tmp/helper"
	"tmp/model"
)

var (
	usercountmutex sync.RWMutex
	usercount      int64
)

// 用户注册
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

// 用户上传订单，将用户添加到缓存池中
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
	atoi, _ := strconv.Atoi(userid)
	_, b := global.GlobalUser.Get("user" + userid)
	if b {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "您有订单未完成",
		})
		return
	}
	order := &model.Order{
		UserId:     int64(atoi),
		Money:      floatmoney,
		CreateTime: createtime,
		StartPlace: startplace,
		EndPlace:   endplace,
		Status:     status,
	}
	err = dao.DB.Create(&order).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订单存储失败" + err.Error(),
		})
		return
	}
	bytes := helper.OrderStruct2Bytes(*order)
	err = helper.Publish("order", bytes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发布消息失败" + err.Error(),
		})
		return
	}
	var curuser model.User
	err = dao.DB.Where("id=?", atoi).Take(&curuser).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户信息出错:",
		})
		return
	}
	// 将用户添加到缓存池中
	err = global.GlobalUser.Add("user"+userid, curuser, 0)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户存储出错:" + err.Error(),
		})
		return
	}
	// 判断是否需要等待
	//ratio := getUserDriverRatio()
	//if ratio > 1 {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"msg":  "预计等待时间:" + string(getUserDriverDiff()) + "分钟",
	//	})
	//	return
	//}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  order,
	})
}

// 司机用户比值
func getUserDriverRatio() float32 {
	return float32(global.GlobalUser.ItemCount() / global.GlobalDriver.ItemCount())
}

// 司机用户差值
func getUserDriverDiff() int {
	return global.GlobalUser.ItemCount() - global.GlobalDriver.ItemCount()
}

// 用户查询订单状态
func UserQueryOrder(c *gin.Context) {
	//userid := c.PostForm("userid")
	orderid := c.PostForm("orderid")
	// 初始版本：查询redis
	//result, err := dao.Rdb.Get(context.Background(), "user"+userid).Result()
	//if err == redis.Nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"msg":  "等待接单中",
	//	})
	//	return
	//}
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "查询订单状态出错",
	//	})
	//	return
	//}
	//c.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//	"msg":  result,
	//})

	// 版本2：使用go-cache,查询订单的缓存对象池
	orderkey := "order" + orderid
	order, find := global.GlobalOrder.Get(orderkey)
	if find {
		m, ok := order.(model.Order)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  m,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "订单信息转换失败",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "订单未找到",
		})
	}
}
