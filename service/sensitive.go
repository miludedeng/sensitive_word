package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var dic map[string]interface{} = make(map[string]interface{})

//初始化词典
func init() {
	dicPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if strings.HasSuffix(dicPath, "_test") {
		_, dicPath, _, _ = runtime.Caller(0)
		dicPath, _ = filepath.Abs(filepath.Dir(dicPath) + "/../")
	}
	//	println(dicPath)
	f, err := os.Open(dicPath + "/dic.txt")
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
	fmt.Println(dic)
}

//添加敏感词到内存
func addWord(word string) {
	tempMap := dic
	runes := []rune(word)
	for i, r := range runes {
		char := string(r)
		if _, ok := tempMap[char].(string); tempMap[char] != nil && !ok {
			tempMap = tempMap[char].(map[string]interface{})
		} else {
			tempMap2 := make(map[string]interface{})
			tempMap2["ISEND"] = "NOEND"
			tempMap[char] = tempMap2
			tempMap = tempMap2
		}
		if len(runes) <= i+1 {
			tempMap["ISEND"] = "END"
		}
	}
}

func SensitiveFind(runes []rune) []string {
	step := 0
	flag := true
	var find func(dictionary map[string]interface{}, key string, result []string, level int) []string
	find = func(dictionary map[string]interface{}, key string, result []string, level int) []string {
		for step < len(runes) {
			char := string(runes[step])
			step++
			if dictionary[char] != nil {
				tempM := dictionary[char].(map[string]interface{})
				result = find(tempM, key+char, result, level+1)
				if level != 0 {
					if tempM["ISEND"] == "END" {
						if flag {
							flag = false
							result = append(result, key+char)
						}
					}
					return result
				}
			} else {
				if level != 0 {
					step--
					return result
				}
			}
			flag = true
		}
		return result
	}
	resultString := find(dic, "", nil, 0)
	//	fmt.Println(resultString)
	return resultString
}

//func SensitiveFind(runes []rune) []string {
//	keywords := []string{}
//	for i := 0; i < len(runes); i++ {
//		nowMap := dic
//		length := 0
//		flaglength := 0
//		flag := false
//		for j := i; j < len(runes);j++ {
//			word := string(runes[j : j+1])
//			nowMap, _ = nowMap[word].(map[string]interface{})
//			if nowMap != nil {
//				length = length + 1
//				tag, _ := nowMap["ISEND"].(string)
//				if "END" == tag {
//					flag = true
//					flaglength = length
//				}
//			} else {
//				break
//			}
//		}
//
//		if length < 2 || !flag {
//			length = 0
//		}
//		if length > 0 {
//			keywords = append(keywords, string(runes[i:i+flaglength]))
//			i = i + flaglength - 1
//		}
//	}
//	fmt.Println(keywords)
//	return keywords
//
//}

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
				switch dictionary[char].(type) {
				case string:
					wordTemp = "**"
				case map[string]interface{}:
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
		if len(runes) == step && wordTemp != "" && level == 0 {
			temp += wordTemp
		}
	}
	find(dic, 0)
	return temp
}

func DoCheck(content string) (words []string, resultContent string) {
	ch := make(chan interface{}, 2)
	// 开启多线程
	go func() {
		ch <- SensitiveFind([]rune(content))
	}()
	go func() {
		ch <- SensitiveReplace([]rune(content))
	}()
	// 从channel中读取数据
	msg := <-ch
	switch msg.(type) {
	case string:
		resultContent, words = (msg).(string), (<-ch).([]string)
	case []string:
		resultContent, words = (<-ch).(string), (msg).([]string)
	}
	return
}
