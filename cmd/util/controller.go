package util

import (
	"errors"
	"fmt"
)

// Controller Controller 数据
type Controller struct {
	Dir string // 存放Controller文件目录名称，默认是根目录下的 router目录
	Pkg string // 接口的版本，如 v1, adm, 同时也是package name
	Tpl string // Controller 模板

	HelperPkg string // 在 /Pkg/下的helper的目录名
	HelperTpl string // Helper 模板

	Model *Model
}

func (c *Controller) FilePath() (string, string, error) {
	src, err := resolveTPLPath(c.Tpl)
	if err != nil {
		return "", "", err
	}
	if c.Model == nil {
		return "", "", errors.New("unknown model")
	}

	dest := fmt.Sprintf("%s/%s/%ss.go", c.Dir, c.Pkg, c.Model.LowerName)
	return src, dest, nil
}

func (c *Controller) HelperFilePath() (string, string, error) {
	src, err := resolveTPLPath(c.HelperTpl)
	if err != nil {
		return "", "", err
	}
	if c.Model == nil {
		return "", "", errors.New("unknown model")
	}

	dest := fmt.Sprintf("%s/%s/%s/%s.go", c.Dir, c.Pkg, c.HelperPkg, c.Model.LowerName)
	return src, dest, nil
}

// ValidHelper 是否需要创建helper文件
func (c *Controller) ValidHelper() bool {
	return c.HelperPkg != "" && c.HelperTpl != ""
}
