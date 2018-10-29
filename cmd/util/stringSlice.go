package util

import "fmt"

// StringSlice 定义字符串数组参数
type StringSlice []string

// String implement flag.Value String method
func (s *StringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

// Set implement flag.Value Set method
func (s *StringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}
