package router

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"igen/lib/logger"
	gMiddleware "igen/lib/middleware"
	"igen/lib/util/httpUtil"
)

func redirect2dfs(format ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		httpUtil.GinRedirect(c, httpUtil.ServiceDFS, format...)
	}
}

func proxy(f func(*gin.Context, string, string, map[string]interface{}), format ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		var path = c.Request.URL.RequestURI()
		var method = strings.ToUpper(c.Request.Method)
		if len(format) > 1 {
			path = strings.Join(format, "")
			if len(c.Params) > 0 {
				ps := []interface{}{}
				for _, v := range c.Params {
					ps = append(ps, v.Value)
				}
				path = fmt.Sprintf(path, ps...)
			}
		}

		var params map[string]interface{}
		if len(format) > 1 || method == "POST" || method == "PUT" || method == "PATCH" {
			var err error
			params, err = gMiddleware.GetParams(c)
			if err != nil {
				logger.Ctx(c).Error("proxy get params error", logger.Err(err))
				c.JSON(200, map[string]interface{}{
					"code": 500,
					"msg":  err.Error(),
				})
				return
			}
		}
		logger.Ctx(c).Debugf("%s %s", method, path)
		f(c, method, path, params)
	}
}
