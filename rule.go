package validate

import (
	"regexp"
)

type Rule struct {
	RuleMethod map[string]interface{}
}

// 初始化规则
func (r *Rule) initRule() *Rule {
	r.RuleMethod = make(map[string]interface{})
	r.RuleMethod["number"] = r.IsNumber
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

// 判断是否在数组中
func (Rule) InArrayString(content string, values []string) bool {
	for _, value := range values {
		if value == content {
			return true
		}
	}
	return false
}
