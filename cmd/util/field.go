package util

// Field 定义model的字段
type Field struct {
	Name    string // 大写开头的字段名 UserID
	Type    string
	TagName string // camel风格  比如userId
	Comment string
}
