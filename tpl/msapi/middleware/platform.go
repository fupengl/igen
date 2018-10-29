package middleware

import (
	"github.com/gin-gonic/gin"
)

// GetPlatform 获取请求来源平台, iOS/Android/Weixin/web/wap/wx
// 用于同一账号多平台同时登录
// TODO
func GetPlatform(c *gin.Context) string {
	return ""
}
