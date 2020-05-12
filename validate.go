package validate

import (
	"reflect"
	"strings"
)

type Validates map[string]string // 验证器所用的数据类型

type Validate struct {
	Rule                      //规则模块
	defaultErrorMsg Validates // 内置默认错误提示
	typeMsg         Validates // 错误规则提示

	rule      map[string][]string // 当前的验证规则
	msg       Validates           // 用户自设提示消息
	scene     map[string][]string // 验证场景
	sceneName string              // 当前验证场景

	checkRule map[string][]string //要检查的规则
	error     string              // 错误消息
}

// 初始化验证类
func (v *Validate) Init() *Validate {
	v.typeMsg = make(Validates)
	v.rule = make(map[string][]string)
	v.msg = make(Validates)
	v.scene = make(map[string][]string)
	v.checkRule = make(map[string][]string)
	v.defaultErrorMsg = make(Validates)
	v.initRule()
	return v
}

// 规则验证
func (v *Validate) Check(values map[string]interface{}, fail bool) bool {
	var (
		ruleValue []string
	)
	for field, rules := range v.checkRule {
		if v.InArrayString("require", rules) || !isEmpty(field, values) { // 判断有必填验证或者值不为空
			for _, rule := range rules {
				rule, ruleValue = getRuleValue(rule)
				if rule == "require" {
					if isEmpty(field, values) {
						v.setError(field, rule)
						return false
					}
				} else {
					value, _ := values[field]
					if !v.callRule(rule, value, rule, values, ruleValue...) {
						v.setError(field, rule, ruleValue)
						if fail {
							panic(v.GetError())
						}
						return false
					}
				}
			}
		}
	}
	return true
}

// 获取自定义规则参数
func getRuleValue(rule string) (newRule string, ruleValue []string) {
	r := strings.Split(rule, ":")
	rLen := len(r)
	if rLen == 1 { // 没有后置参数
		return rule, nil
	} else if rLen == 2 {
		param := strings.Split(r[1], ",")
		return r[0], param
	} else { // 多个：无法处理后续参数，只返回规则rule
		return r[0], nil
	}
}

// 是否为空值
func isEmpty(field string, values map[string]interface{}) bool {
	value, ok := values[field]
	if !ok {
		return true
	}
	t := reflect.TypeOf(value)
	if t.Name() == "string" {
		if value == "" {
			return true
		}
	}
	return false
}

// 设置验证场景
func (v *Validate) Scene(sceneName string) *Validate {
	v.sceneName = sceneName
	if v.sceneName == "" {
		v.checkRule = v.rule
	} else {
		ruleKey, ok := v.scene[sceneName]
		if !ok {
			// TODO:自定义场景函数
			v.checkRule = v.rule
			return v
		}

		checkRule := make(map[string][]string)
		for _, value := range ruleKey {
			ruleValue, rvOK := v.rule[value]
			if rvOK {
				checkRule[value] = ruleValue
			}
		}
		v.checkRule = checkRule
	}
	return v
}

// 设置验证场景数据
func (v *Validate) SetScene(scene Validates) *Validate {
	// 直接赋值会覆盖掉，这样写,重复的覆盖，不重复的新增
	var sceneKey []string
	for key, value := range scene {
		sceneKey = strings.Split(value, ",")
		v.scene[key] = sceneKey
	}
	return v
}

// 设置规则
func (v *Validate) SetRule(rule Validates) *Validate {
	var ruleArray []string
	for key, value := range rule {
		ruleArray = strings.Split(value, "|")
		v.rule[key] = ruleArray
	}
	v.Scene(v.sceneName)
	return v
}

// 设置消息规则
func (v *Validate) SetMsg(msg Validates) *Validate {
	// 直接赋值会覆盖掉，这样写,重复的覆盖，不重复的新增
	for key, value := range msg {
		v.msg[key] = value
	}
	return v
}

// 设置错误信息
func (v *Validate) setError(keyName string, ruleName string, ruleValue ...interface{}) {
	errorMsg, ok := v.msg[keyName+"."+ruleName]
	if ok {
		v.error = errorMsg
	} else {
		errorMsg, ok = v.typeMsg[ruleName]
		if ok {
			v.error = errorMsg
		} else {
			v.error = keyName + " rule[" + ruleName + "] validation failed"
		}
	}
}

// 读取错误信息
func (v *Validate) GetError() string {
	return v.error
}
