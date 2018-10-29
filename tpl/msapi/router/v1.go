package router

import (
	"github.com/gin-gonic/gin"

	"igen/msdemo/router/v1"
)

func v1RouterNoSignNoToken(rg *gin.RouterGroup) {
	rg.POST("/pictures", redirect2dfs("/v1/pictures"))
}

// 不需要验证Sign的接口
func v1RouterNoSign(rg *gin.RouterGroup) {

}

// 不需要验证Token的接口
func v1RouterNoToken(rg *gin.RouterGroup) {
	// 手机验证码
	rg.POST("/sms", v1.CreateSMS)

	// 登录
	rg.POST("/sessions", v1.CreateSession)
}

// 需要验证Sign和Token的接口
func v1Router(rg *gin.RouterGroup) {
	// 登出
	rg.DELETE("/sessions", v1.DeleteSession)
}
