package service

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"os"
	"strings"
)

var dic map[string]interface{} = make(map[string]interface{})

//初始化词典
func init() {
	f, err := os.Open(beego.AppPath + "/dic.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		addWord(strings.Replace(line, "\n", "", -1)) //文件的没一行有\n结尾，需要剔除
	}

}

//添加敏感词到内存
func addWord(word string) {
	tempMap := dic
	runes := []rune(word)
	for i, r := range runes {
		if len(runes) <= i+1 {
			break
		}
		char := string(r)
		nChar := string(runes[i+1])
		if tempMap[char] != nil && fmt.Sprintf("%T", tempMap[char]) != "string" {
			tempMap = tempMap[char].(map[string]interface{})
			if tempMap[nChar] == nil {
				tempMap[nChar] = "END"
			}

		} else {
			tempMap2 := make(map[string]interface{})
			tempMap2[nChar] = "END"
			tempMap[char] = tempMap2
			tempMap = tempMap2
		}
	}
}

//敏感词查询
func SensitiveFind(runes []rune) []string {
	step := 0
	var find func(dictionary map[string]interface{}, key string, result []string, level int) []string
	find = func(dictionary map[string]interface{}, key string, result []string, level int) []string {
		for step < len(runes) {
			char := string(runes[step])
			step++
			if dictionary[char] != nil {
				if fmt.Sprintf("%T", dictionary[char]) == "string" {
					result = append(result, key+char)
					return result
				}
				result = find(dictionary[char].(map[string]interface{}), key+char, result, level+1)
				if level != 0 {
					return result
				}
			} else {
				if level != 0 {
					step--
					return result
				}
			}
		}
		return result
	}
	return find(dic, "", nil, 0)
}

//敏感词替换
func SensitiveReplace(runes []rune) string {
	step := 0
	var temp string
	var wordTemp string
	var find func(dictionary map[string]interface{}, level int)
	find = func(dictionary map[string]interface{}, level int) {
		for step < len(runes) {
			char := string(runes[step])
			step++
			if dictionary[char] != nil {
				wordTemp += char
				if fmt.Sprintf("%T", dictionary[char]) == "string" {
					wordTemp = "**"
					if step == len(runes) {
						temp += wordTemp
					}
				} else {
					find(dictionary[char].(map[string]interface{}), level+1)
				}
				if level != 0 {
					return
				}
			} else {
				if level == 0 {
					if wordTemp != "" {
						temp += wordTemp
						wordTemp = ""
					}
					temp += char
				} else {
					step--
					return
				}
			}
		}
	}
	find(dic, 0)
	return temp
}
