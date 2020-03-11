# validate
Go语言验证类

### 内置规则

| 规则 | 备注 | 示例 |
|:-|:-|:-|
|number|判断字段是否为数字|number|
|chsAlohaNum|判断字段是否由汉字、字母和数字组成|chsAlohaNum|
|alphaNum|判断字段是否由字母和数字组成|alphaNum|
|length|判断字段长度是否符合要求|length:1,2|
|colorHex|判断字段是否为16进制的颜色值|colorHex|
|chsDash|验证某个字段的值只能是汉字、字母、数字和下划线_及破折号-|chsDash|
|chsAlpha|验证某个字段的值只能是汉字、字母|chsAlpha|
|chs|验证某个字段的值只能是汉字|chs|
|between|验证某个字段的值是否在某个区间(数值类型)|between:1,10|
|notBetween|验证某个字段的值不在某个范围(数值类型)|notBetween:1,10|
|in|验证某个字段的值是否在指定的值内|in:1,2,3,4|
|notIn|验证某个字段的值是否不在指定的值内|notIn:1,2,3,4|
|max|最大值限制|max:10|
|min|最小值限制|min:10|