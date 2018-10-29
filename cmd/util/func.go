package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// toLowerCamelCase 把大写的字段名变成小写的 同时处理特殊 URL ID
// example:
// IDCard --> idCard
// ParentID --> parentId
// ImageURL --> imageUrl
func toLowerCamelCase(str string) string {
	//把word划分为单词 同时根据单词的位置及内容变换单词
	var words []string
	var lastWordBegin = 0

	end := len(str) - 1

	for cursor, char := range str {

		//已经到最后一个字符
		if cursor == end {
			words = appendWord(words, str[lastWordBegin:cursor+1])
			break
		}

		//小写 后面接着大写  example: aC
		if char >= 'a' && char <= 'z' && str[cursor+1] >= 'A' && str[cursor+1] <= 'Z' {
			words = appendWord(words, str[lastWordBegin:cursor+1])
			lastWordBegin = cursor + 1
			continue
		}

		//大写 后面接着小写 前面接着大写  example: DCa
		if cursor != 0 && char >= 'A' && char <= 'Z' &&
			str[cursor+1] >= 'a' && str[cursor+1] <= 'z' &&
			str[cursor-1] >= 'A' && str[cursor-1] <= 'Z' {
			words = appendWord(words, str[lastWordBegin:cursor])
			lastWordBegin = cursor
			continue
		}

	}

	return strings.Join(words, "")
}

// ToUnderScore
// eg: UserName -> user_name
func ToUnderScore(str string) string {
	char := []int32{}
	for _, s := range str {
		if s >= 'A' && s <= 'Z' {
			char = append(char, '_')
			char = append(char, s+32)
		} else {
			char = append(char, s)
		}
	}
	return string(char)
}

func appendWord(words []string, word string) []string {
	first := len(words) == 0
	if first {
		words = append(words, strings.ToLower(word))
	} else {
		//这两个单词 特殊处理 输出为Id Url
		if word == "ID" || word == "URL" {
			words = append(words, word[:1]+strings.ToLower(word[1:]))
		} else {
			words = append(words, word)
		}
	}

	return words
}

func resolveTPLPath(f string) (string, error) {
	if f == "" {
		return f, errors.New("tpl is empty")
	}
	//首先在当前目录查找
	if _, err := os.Stat(f); err == nil {
		return f, nil
	}
	//接着在../igen查找
	goPath := os.Getenv("GOPATH")
	tplPath := goPath + "/src/igen/tpl/" + f

	if _, err := os.Stat(tplPath); err == nil {
		return tplPath, nil
	}
	return "", fmt.Errorf("can not find %s in: %s \n ", f, tplPath)
}

func getProjectName() string {
	p, _ := os.Getwd()

	if i := strings.LastIndex(p, "/"); i >= 0 {
		p = p[0:i]
	}

	if i := strings.LastIndex(p, "/"); i >= 0 {
		p = p[i+1:]
	}
	return p
}

func getSubProjectName() string {
	pwd, _ := os.Getwd()
	return path.Base(pwd)
}

func Gofmt(filePath string) {
	cmd := exec.Command("gofmt", "-w", filePath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("gofmt err '%s' , please run : gofmt -w %s\n", err.Error(), filePath)
	}
}
