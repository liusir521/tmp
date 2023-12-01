package service

import (
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"net/http"
	"strconv"
	"time"
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
		return
	}
	// 连接消息队列
	err = consumer.ConnectToNSQD(conf.NsqdAddr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "连接消息队列出错:" + err.Error(),
		})
		return
	}
	<-consumer.StopChan
}

// 完成订单（普通订单）
func FinshOrder(c *gin.Context) {
	userid := "user" + c.PostForm("userid")
	orderid := "order" + c.PostForm("orderid")
	driverid := "driver" + c.PostForm("driverid")
	// 将用户从其缓存池中删除
	userres, b := global.GlobalUser.Get(userid)
	if b {
		user := userres.(model.User)
		get, b := global.GlobalOrder.Get(orderid)
		if b {
			order := get.(model.Order)
			// 扣除用户费用，将用户从缓存池中删除
			user.Money -= order.Money
			err2 := dao.DB.Save(&user).Error
			if err2 != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "用户更新失败" + err2.Error(),
				})
				return
			}
			get, b := global.GlobalDriver.Get(driverid)
			if b {
				// 更新司机信息
				driver := get.(model.Driver)
				driver.Money += order.Money
				driver.RunCount++
				err2 := dao.DB.Save(&driver).Error
				if err2 != nil {
					c.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg":  "司机更新失败" + err2.Error(),
					})
					return
				}
				global.GlobalDriver.Set(driverid, driver, 0)
				order.Status = "已完成"
				// 将订单信息更新到数据库
				err := dao.DB.Save(&order).Error
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg":  "订单更新失败" + err.Error(),
					})
					return
				}
				// 将订单从订单缓存池中删除
				global.GlobalOrder.Delete(orderid)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "司机信息出错",
				})
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "订单信息出错",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户信息出错",
		})
		return
	}
}

// 发布顺风车订单，让用户消费
func TogetherOrderUpload(c *gin.Context) {
	driverid := "driver" + c.PostForm("driverid")
	counti := c.PostForm("count") // 搭乘人数
	count, _ := strconv.Atoi(counti)
	// 判断司机是否已下线
	driverinc, b := global.GlobalDriver.Get(driverid)
	if !b {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "司机状态有误",
		})
	}
	driver := driverinc.(model.Driver)

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
		DriverId:   driver.ID,
		Carid:      driver.CarId,
		Money:      floatmoney,
		CreateTime: createtime,
		StartPlace: startplace,
		EndPlace:   endplace,
		Status:     status,
	}

	// 发布消息
	bytes := helper.OrderStruct2Bytes(*order)
	for i := 0; i < count; i++ {
		err = helper.Publish(startplace+endplace, bytes)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "发布消息失败" + err.Error(),
			})
			return
		}
	}
	// 向缓冲池中添加新的订单，存储订单信息
	global.GlobalOrder.Add("driver"+string(driver.ID), order, 0)
	// GlobalTogetherOrder缓冲池来代表司机是否发车
	err = global.GlobalTogetherOrder.Add("driver"+string(driver.ID), 0, 0)
	// 向sync.map中初始化key和value，用于存放当前顺风车的乘客ID
	userids := []int64{}
	global.TogetherUsers.Store("driver"+string(driver.ID), userids)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "司机存储失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "订单发布成功",
	})
}

// 顺风车司机发车
func StartTogetherOrder(c *gin.Context) {
	driverid := "driver" + c.PostForm("driverid")
	global.GlobalTogetherOrder.Delete(driverid)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "司机发车成功",
	})
}

func FinshTogetherOrder(c *gin.Context) {
	driverid := "driver" + c.PostForm("driverid")
	get, b := global.GlobalOrder.Get(driverid)
	if b {
		order := get.(model.Order)
		useridsinc, isfind := global.TogetherUsers.Load(driverid)
		if !isfind {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "司机信息有误",
			})
			return
		}
		userids := useridsinc.([]int64)
		drivermoney, usermoney := order.Money, order.Money/float64(len(userids))
		for _, userid := range userids {
			order.UserId = userid
			order.Money = usermoney
			err := dao.DB.Create(&order).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "用户订单存储有误",
				})
			}
		}
		driverorder := &model.TgOrder{
			UserCouunt: int64(len(userids)),
			DriverId:   order.DriverId,
			Carid:      order.Carid,
			Money:      drivermoney,
			CreateTime: order.CreateTime,
			StartPlace: order.StartPlace,
			EndPlace:   order.EndPlace,
			Status:     "已完成",
		}
		dao.DB.Create(&driverorder)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "订单完成",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "订单信息有误",
		})
	}
}
