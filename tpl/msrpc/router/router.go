package router

import (
	"github.com/gin-gonic/gin"

	"igen/lib/consul"
	"igen/lib/render"
)

// Init routers
func Init(r *gin.Engine) {

	// 404 Not Found
	r.NoRoute(func(c *gin.Context) {
		render.Err404(c)
	})

	consulRouter(r.Group("/consul"))
	v1Router(r.Group("/v1"))
}

// Consul health check
func consulRouter(rg *gin.RouterGroup) {

	rg.GET("/actions/check_http", func(c *gin.Context) {
		c.Set("isOk", true)
		consul.Check(c.Writer, c.Request)
	})

	rg.GET("/actions/check_rpc", func(c *gin.Context) {
		c.Set("isOk", true)
		consul.Check(c.Writer, c.Request)
	})

}
