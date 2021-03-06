package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"sensitive_word/service"
	"sensitive_word/util"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Sensitive() {
	content := c.GetString("content")
	begin := time.Now()
	words, resultContent := service.DoCheck(content)
	since := time.Since(begin).Nanoseconds() //统计查找敏感词和替换敏感词所用的时间
	result := &struct {
		Words   map[string]int `json:"words"`
		Count   int            `json:"count"`
		Content string         `json:"content"`
		Cost    string         `json:"cost"`
	}{
		Words:   util.SliceDuplicClear(words),
		Count:   len(words),
		Content: resultContent,
		Cost:    fmt.Sprintf("%f%s", float64(since)/1000000, "ms"),
	}
	c.Data["json"] = &result //返回json类型数据
	c.ServeJSON()
}
