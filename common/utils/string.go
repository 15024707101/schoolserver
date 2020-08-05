package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// UnderscoreName 驼峰式写法转为下划线写法
func UnderscoreName(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}

	return buffer.String()
}

// CamelName 下划线写法转为驼峰写法
func CamelName(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// SearchString 搜索字符串
func SearchString(slice []string, s string) int {
	for i, v := range slice {
		if s == v {
			return i
		}
	}

	return -1
}

// UcFirst 首字母大写
func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}

	return ""
}

func GeneratePassword() string {
	res := ""
	for i := 0; i < 2; i++ {
		res += RandomString(4)
		res += "-"
	}
	res += RandomString(4)
	return res
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLKMNOPQRSTUVWXYZ"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//SubStr 截取字符串
func SubStr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// float64 转字符串
func ValueOfFloat64(value float64) string {
	return strconv.FormatFloat(float64(value), 'f', 6, 64)
}

// int64 转字符串
func ValueOfInt(value int64) string {
	return strconv.FormatInt(value, 10)
}

// IsInSlice
// 字符串是否在slice中
func IsInSlice(slice []string, s string) bool {
	for _, thisS := range slice {
		if thisS == s {
			return true
		}
	}
	return false
}

// AtoIDefault0
// 字符串转int,不成功则0
func AtoIDefault0(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// FormatInt
// 格式化int为string
func FormatInt(i int) string {
	return strconv.Itoa(i)
}

// FormatFloat
// 格式化float64为string
func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// FormatFloatPrec0
// 格式化float64位字符串，只取整数位
func FormatFloatPrec0(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64)
}

// FormatFloatPrec0
// 格式化float64位字符串，2位小数点
func FormatFloatPrec2(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// FormatFloatPrec0
// 格式化float64位字符串,4位小数点
func FormatFloatPrec4(f float64) string {
	return strconv.FormatFloat(f, 'f', 4, 64)
}

// ParseFloat64
// 将string解析为float64
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseFloat64Default0
// 将string解析为float64，失败为0
func ParseFloat64Default0(s string) float64 {
	out, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return out
}

// IntArrayToStringArr
// int数组转换为string数组
func IntArrayToStringArr(i []int) []string {
	strArr := make([]string, len(i))
	for k, v := range i {
		strArr[k] = FormatInt(v)
	}
	return strArr
}
