package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"tmp/dao"
	"tmp/helper"
	"tmp/model"
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

func UserUploadOrder(c *gin.Context) {
	//userid := c.PostForm("userid")
	//money := c.PostForm("money")
	//floatmoney, err := strconv.ParseFloat(money, 64)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "金额格式错误",
	//	})
	//	return
	//}
	//createtime := time.Now()
	//startplace := c.PostForm("startplace")
	//endplace := c.PostForm("endplace")
	//status := "已创建"
	//order := &model.Order{
	//	UserId:     userid,
	//	Money:      floatmoney,
	//	CreateTime: createtime,
	//	StartPlace: startplace,
	//	EndPlace:   endplace,
	//	Status:     status,
	//}
	//bytes := helper.OrderStruct2Bytes(*order)
	//fmt.Println(string(bytes))
	order2 := &model.Order{
		UserId:     "1",
		Money:      1,
		CreateTime: time.Now(),
		StartPlace: "北京",
		EndPlace:   "上海",
		Status:     "已创建",
	}
	bytes2 := helper.OrderStruct2Bytes(*order2)
	var err error
	err = helper.Publish("order", bytes2)
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
