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
	r.RuleMethod["notBetween"] = r.NotBetween
	r.RuleMethod["between"] = r.Between
	r.RuleMethod["chs"] = r.IsChs
	r.RuleMethod["chsAlpha"] = r.IsChsAlpha
	r.RuleMethod["chsDash"] = r.IsChsDash
	r.RuleMethod["in"] = r.In
	r.RuleMethod["notIn"] = r.NotIn
	r.RuleMethod["max"] = r.Max
	r.RuleMethod["min"] = r.Min
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
		panic("rule " + name + " not fount")
	}
}

// 判断字段是否为数字
func (Rule) IsNumber(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	switch value.(type) {
	case string:
		result, _ := regexp.MatchString("\\d+", value.(string))
		return result
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	}
	return false
}

// 判断字段是否由汉字、字母和数字组成
func (Rule) IsChsAlphaNum(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
	m, _ := regexp.MatchString("^[\u4e00-\u9fa5a-zA-Z0-9]+$", value.(string))
	return m
}

// 判断字段是否由字母和数字组成
func (Rule) IsAlphaNum(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
	m, _ := regexp.MatchString("^[A-Za-z0-9]+$", value.(string))
	return m
}

// 判断字段长度是否符合要求
func (Rule) Length(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
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
	value = String(value)
	m, _ := regexp.MatchString("^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$", value.(string))
	return m
}

// 验证某个字段的值只能是汉字、字母、数字和下划线_及破折号-
func (Rule) IsChsDash(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
	m, _ := regexp.MatchString("^[\u4e00-\u9fa5a-zA-Z0-9_\\-]+$", value.(string))
	return m
}

// 验证某个字段的值只能是汉字、字母
func (Rule) IsChsAlpha(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
	m, _ := regexp.MatchString("^[\u4e00-\u9fa5a-zA-Z]+$", value.(string))
	return m
}

// 验证某个字段的值只能是汉字
func (Rule) IsChs(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
	m, _ := regexp.MatchString("^[\u4e00-\u9fa5]+$", value.(string))
	return m
}

// 验证某个字段的值是否在某个区间(数值类型)
func (r *Rule) Between(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	if len(arg) != 2 {
		panic("Rule BetWeen Error")
		return false
	}
	value = Float64(value)
	min := Float64(arg[0])
	max := Float64(arg[1])
	if value.(float64) > max || value.(float64) < min {
		return false
	}

	return true
}

// 验证某个字段的值不在某个范围(数值类型)
func (r *Rule) NotBetween(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	if len(arg) != 2 {
		panic("Rule NotBetween Error")
		return false
	}
	if !r.IsNumber(value, rule, nil) {
		return false
	}
	return !r.Between(value, rule, data, arg...)
}

// 验证某个字段的值是否在指定的值内
func (r *Rule) In(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	value = String(value)
	return r.InArrayString(value.(string), arg)
}

// 验证某个字段的值是否不在指定的值内
func (r *Rule) NotIn(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	return !r.In(value, rule, data, arg...)
}

/**
 * 最大值限制
 * 当类型为数值时，判断数值大小
 * 当类型为切片/字符串/通道时，判断成员数或者文本长度
 * 其他类型不进行判断
 */
func (r *Rule) Max(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	if len(arg) < 1 {
		panic("Rule Max Error")
		return false
	}
	if r.IsNumber(value, rule, data) {
		return Float64(value) <= Float64(arg[0])
	}
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	switch kind {
	case reflect.String:
		return len(value.(string)) <= Int(arg[0])
	}
	return true
}

/**
 * 最小值限制
 * 当类型为数值时，判断数值大小
 * 当类型为切片/字符串/通道时，判断成员数或者文本长度
 * 其他类型不进行判断
 */
func (r *Rule) Min(value interface{}, rule string, data map[string]interface{}, arg ...string) bool {
	if len(arg) < 1 {
		panic("Rule Max Error")
		return false
	}
	if r.IsNumber(value, rule, data) {
		return Float64(value) > Float64(arg[0])
	}
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	switch kind {
	case reflect.String:
		return len(value.(string)) > Int(arg[0])
	}
	return true
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
