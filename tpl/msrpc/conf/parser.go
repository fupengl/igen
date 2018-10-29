package conf

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"igen/lib/consul"
)

// parseToml 把toml模板值替换为环境变量的值
func parseToml(filepath string) (io.Reader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer

	kv := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := parseLine(scanner.Text(), kv)
		buffer.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buffer.Bytes()), nil
}

// 解析toml文件每行{{def|env}}格式的字符串
// 优先从env环境变量中取值，如果env中没有，则取默认的def
//
// 还支持{{$var}}格式， 如果上文中有var的值，$var将替换为那个值
//
func parseLine(str string, kv map[string]string) string {
	l := len(str)
	if l < 5 {
		return str
	}

	out := []byte{}
	mark := -1

	var (
		match, split   bool
		defVal, relVal []byte
		i, equalIndex  int
	)

	for i = 0; i < l; i++ {
		// {{ 开始
		if str[i] == 123 && str[i+1] == 123 {
			if mark > -1 {
				out = out[0:mark]
				out = append(out, str[mark:i]...)
			}

			match = true
			split = false
			mark = i

			defVal = []byte{}
			relVal = []byte{}

			i++
			continue
		}

		if match {

			// }} 结束
			if str[i] == 125 && str[i+1] == 125 {
				i++
				match = false
				split = false
				mark = -1

				out = append(out, fixVal(defVal, relVal, kv)...)
				continue
			}

			// 最后 没有闭合的 }}
			if i == l-1 {
				out = out[0:mark]
				out = append(out, str[mark:]...)
			}

			// 空格， 剔除空格
			if str[i] == 32 {
				continue
			}

			// | 竖线
			if str[i] == 124 {
				split = true
				continue
			}

			if split {
				relVal = append(relVal, str[i])
			} else {
				defVal = append(defVal, str[i])
			}

			continue
		}

		out = append(out, str[i])

		// = 第一个等号
		if equalIndex == 0 && str[i] == 61 && i+1 < l {
			equalIndex = i
		}
	}

	// 有等号就当作是 k-v
	if equalIndex > 0 {
		val := strings.TrimSpace(string(out[equalIndex+1:]))
		l := len(val)
		if l > 1 && val[0] == val[l-1] && val[0] == 34 { // "
			// 去除双引号
			kv[strings.TrimSpace(string(out[0:equalIndex]))] = string(val[1 : l-1])
		} else {
			kv[strings.TrimSpace(string(out[0:equalIndex]))] = val
		}
	}

	return string(out)
}

func fixVal(defVal, relVal []byte, kv map[string]string) []byte {

	if len(relVal) > 0 {
		rv := string(relVal)

		// 优先取配置中心的数据
		if v, err := consul.Get(rv); err == nil {
			log.Printf("configuation center: %s=%s", rv, string(v))
			return v
		}

		// 再取环境变量
		if v := os.Getenv(rv); v != "" {
			log.Printf("env        : %s=%s", rv, v)
			return []byte(v)
		}
	}

	// $变量
	if len(defVal) > 0 && defVal[0] == 36 {
		name := defVal[1:]
		if v, exists := kv[string(name)]; exists {
			return []byte(v)
		}
		return name
	}

	return defVal
}
