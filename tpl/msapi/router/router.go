package router

import (
	"github.com/gin-gonic/gin"

	"igen/lib/render"
	"igen/msdemo/middleware"
)

// Init routers
func Init(r *gin.Engine) {

	// 404 Not Found
	r.NoRoute(func(c *gin.Context) {
		render.Err404(c)
	})

	v1Router(r.Group("/v1", middleware.CheckSign(), middleware.CheckToken()))
	v1RouterNoSign(r.Group("/v1", middleware.CheckToken()))
	v1RouterNoToken(r.Group("/v1", middleware.CheckSign()))
	v1RouterNoSignNoToken(r.Group("/v1"))
}
