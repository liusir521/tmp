package router

import (
	"github.com/gin-gonic/gin"
	"tmp/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("user/register", service.UserRegister)
	r.POST("user/uploadorder", service.UserUploadOrder)
	r.POST("driver/register", service.DriverRegister)
	r.POST("driver/work", service.DriverWork)
	r.POST("driver/notwork", service.DriverNotWork)
	return r
}
