package validate

import (
	"reflect"
	"regexp"
	"strconv"
)

type Rule struct {
	RuleMethod map[string]interface{}
}

// 初始化规则
func (r *Rule) initRule() *Rule {
	r.RuleMethod = make(map[string]interface{})
	r.RuleMethod["number"] = r.IsNumber
	r.RuleMethod["chsAlphaNum"] = r.IsChsAlphaNum
	r.RuleMethod["alphaNum"] = r.IsAlphaNum
	r.RuleMethod["length"] = r.Length
	r.RuleMethod["colorHex"] = r.IsColorHex
	return r
}

// 添加自定义规则
func (r *Rule) AddRule(name string, fun func(value interface{}, rule string, data map[string]interface{}, arg ...string) bool) *Rule {
	r.RuleMethod[name] = fun
	return r
}

// 调用规则
func (r *Rule) callRule(name string, value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	fun, ok := r.RuleMethod[name]
	if ok {
		return fun.(func(value interface{}, rule string, data map[string]interface{}, arg ...string) bool)(value, rule, data, arg...)
	} else {
		panic("rule " + name + "not fount")
	}
}

// 判断字段是否为数字，可独立调用
func (Rule) IsNumber(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	switch value.(type) {
	case string:
		result, _ := regexp.MatchString("\\d+", value.(string))
		return result
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

// 判断字段是否由汉字、字母和数字组成
func (Rule) IsChsAlphaNum(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	t := reflect.TypeOf(value)
	if t.Name() != "string" {
		return false
	}
	m, _ := regexp.MatchString("^[\u4e00-\u9fa5a-zA-Z0-9]+$", value.(string))
	return m
}

// 判断字段是否由字母和数字组成
func (Rule) IsAlphaNum(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	t := reflect.TypeOf(value)
	if t.Name() != "string" {
		return false
	}
	m, _ := regexp.MatchString("^[A-Za-z0-9]+$", value.(string))
	return m
}

// 判断字段长度是否符合要求
func (Rule) Length(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	t := reflect.TypeOf(value)
	if t.Name() != "string" {
		return false
	}
	vLen := len(value.(string))
	if len(arg) == 2 { // 判断长度区间
		min, _ := strconv.Atoi(arg[0])
		max, _ := strconv.Atoi(arg[1])
		if vLen > max || vLen < min {
			return false
		}
	} else if len(arg) == 1 { // 判断是否长度一直
		validateLen, _ := strconv.Atoi(arg[0])
		return vLen == validateLen
	} else {
		panic("Rule Length Error")
		return false
	}

	return true
}

// 判断字段是否为16进制的颜色值
func (Rule) IsColorHex(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	t := reflect.TypeOf(value)
	if t.Name() != "string" {
		return false
	}
	m, _ := regexp.MatchString("^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$", value.(string))
	return m
}

// 判断是否在数组中
func (Rule) InArrayString(content string, values []string) bool {
	for _, value := range values {
		if value == content {
			return true
		}
	}
	return false
}
