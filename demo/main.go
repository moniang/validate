package main

import (
	"fmt"
	"github.com/moniang/validate"
)

func main() {
	value := map[string]interface{}{
		"user": 111,
	}
	var v validate.Validate
	v.Init() // 初始化验证类

	v.AddRule("check_user", func(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
		fmt.Println("收到了自定义参数", arg)
		return false
	}) // 添加自定义规则

	v.SetScene(map[string]string{
		"login": "user,pass",
	}) // 设置验证场景字段信息

	v.Scene("login") // 进入验证场景

	v.SetRule(map[string]string{
		"user": "require|number|check_user:1,2,你好",
		"pass": "number",
		"vali": "require",
	}) // 设置验证规则

	v.SetMsg(map[string]string{
		"user.require":    "Value参数必须填写",
		"user.number":     "Value参数必须为数字",
		"user.check_user": "自定义规则错误",
	}) // 设置提示消息

	if v.Check(value) { // 进行判断
		fmt.Println(v.GetError()) // 输出错误信息
	}

}
