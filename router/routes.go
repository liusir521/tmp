package router

import (
	"github.com/gin-gonic/gin"
	"tmp/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	// 用户
	r.POST("user/register", service.UserRegister)
	r.POST("user/uploadorder", service.UserUploadOrder)

	// 司机
	r.POST("driver/register", service.DriverRegister)
	r.POST("driver/work", service.DriverWork)
	r.POST("driver/notwork", service.DriverNotWork)
	r.POST("driver/start", service.DriverStart)
	r.POST("driver/finshorder", service.FinshOrder)
	return r
}
