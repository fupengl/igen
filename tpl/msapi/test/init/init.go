package init

import (
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"

	"igen/msdemo/conf"
)

// 用于测试
func init() {
	gin.SetMode(gin.TestMode)

	_, file, _, _ := runtime.Caller(1)
	path, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../../../")))

	conf.SetEnv(conf.EnvTest)
	conf.AppPath(path)
	conf.Init()
}
