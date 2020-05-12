package main

import (
	"fmt"
	"github.com/moniang/validate"
)

type loginValidate struct {
	validate.Validate
}

func (l *loginValidate) InitValidate() *loginValidate {
	// 初始化验证器
	l.Init()

	// 设置验证规则
	l.SetRule(validate.Validates{
		"user": "require|number",
		"pass": "length:6,20|check_pass",
	})

	// 添加check_pass自定义规则
	l.AddRule("check_pass", func(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
		return (validate.String(value) != "123456") && (validate.Int(data["user"]) != 111)
	})

	// 设置验证错误消息
	l.SetMsg(validate.Validates{
		"user.require":    "Value参数必须填写",
		"user.number":     "Value参数必须为数字",
		"pass.length":     "密码长度为6~20位",
		"pass.check_pass": "账号或密码错误",
	})

	// 设置验证场景字段信息，如不设置，则默认验证所有字段
	l.SetScene(validate.Validates{
		"login": "user,pass",
	})

	// 进入验证场景
	l.Scene("login")

	return l
}

func main() {
	value := map[string]interface{}{
		"user": 111,
		"pass": "123456",
	}

	defer func() {
		fmt.Println("验证异常:", recover())
	}()

	var loginValidate loginValidate
	if loginValidate.InitValidate().Check(value, true) {
		fmt.Println("验证通过")
	} else {
		fmt.Println(loginValidate.GetError())
	}

	// 直接调用使用
	var v validate.Validate
	v.Init() // 初始化验证类
	v.AddRule("check_user", func(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
		fmt.Println("收到了自定义参数", arg)
		return true
	}) // 添加自定义规则

	v.SetScene(validate.Validates{
		"login": "user,pass",
		"test":  "a,b",
	}) // 设置验证场景字段信息

	v.Scene("test") // 进入验证场景

	v.SetRule(validate.Validates{
		"user": "require|number|check_user:1,2,你好",
		"pass": "length:6,20",
		"vali": "require|length:4",
		"a":    "between:3,20",
		"b":    "notIn:1,2,3|min:0",
	}) // 设置验证规则

	v.SetMsg(validate.Validates{
		"user.require":    "Value参数必须填写",
		"user.number":     "Value参数必须为数字",
		"user.check_user": "自定义规则错误",
		"pass.length":     "密码长度为6~20位",
		"vali.require":    "验证码必须填写",
		"vali.length":     "验证码长度错误",
	}) // 设置提示消息

	if !v.Check(value, false) { // 进行判断
		fmt.Println(v.GetError()) // 输出错误信息
	} else {
		fmt.Println("验证通过")
	}

}
