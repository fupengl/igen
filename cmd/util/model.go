package util

import "fmt"

// Model数据
type Model struct {
	Name       string // 资源名称
	LowerName  string
	SName      string // 要操作的model struct的名称，如果为空，则和name值一样
	LowerSName string
	Pkg        string   // model的包名
	Tpl        string   // model模板
	Fields     []*Field // 字段名
	Comment    string   // model描述
}

func (m *Model) FilePath() (string, string, error) {
	src, err := resolveTPLPath(m.Tpl)
	if err != nil {
		return "", "", err
	}
	return src, m.SourcePath(), nil
}

func (m *Model) SourcePath() string {
	return fmt.Sprintf("%s/%s.go", m.Pkg, m.LowerName)
}

func (m *Model) ParseFields() ([]*Field, error) {
	var err error
	m.Fields, err = GetStructFields(m.Name, fmt.Sprintf("%s/%s.go", m.Pkg, m.LowerName))
	return m.Fields, err
}
