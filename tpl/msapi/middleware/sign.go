package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"

	"igen/lib/constant"
	"igen/lib/logger"
	gMiddleware "igen/lib/middleware"
	"igen/lib/render"
	"igen/msdemo/conf"
)

// CheckSign 验证sign
func CheckSign() func(*gin.Context) {
	return func(c *gin.Context) {
		// 非生产环境, 为了方便测试组测试
		if !conf.IsProd() && c.Request.Header.Get(constant.SignNil) == "1" {
			c.Next()
			return
		}

		err := gMiddleware.ValidSign(
			c,
			gMiddleware.FuncAppSecretByAppID(func(appID string) (string, error) {
				return getAppSecretByAppID(appID)
			}),
			"osType",
			"iOS", "ios", "android", "Android", "web", "wap", "weixin", "wx",
		)
		if err != nil {
			logger.Ctx(c).Error(err.Error())
			render.Abort(c, 401, "SIGN授权认证失败")
			return
		}
		c.Next()
	}
}

func getAppSecretByAppID(appID string) (string, error) {
	if v, exists := conf.App.AppSecrets[appID]; exists {
		return v, nil
	}
	return "", errors.New("not found")
}
