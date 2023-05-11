package strUtil

import (
	"encoding/base64"
	"github.com/daida459031925/common/fmt"
	"net"
	"net/url"
	"strings"
	"unicode"
)

// TrimSpaceAndEmpty 将字符串 s 去除首尾空格后，如果为空字符串，则返回 true，否则返回 false。
func TrimSpaceAndEmpty(s string) bool {
	return fmt.Trim(s) == ""
}

// FirstNonEmpty 返回非空字符串数组 arr 中的第一个字符串，如果 arr 中所有字符串都为空字符串，则返回空字符串。
func FirstNonEmpty(arr ...string) string {
	for _, s := range arr {
		if !TrimSpaceAndEmpty(s) {
			return s
		}
	}
	return ""
}

// Concat 将字符串数组 arr 中的所有字符串拼接起来，使用 sep 分隔符分隔。
func Concat(arr []string, sep string) string {
	return strings.Join(arr, sep)
}

// CountOccurrences 统计字符串 s 中子字符串 sub 出现的次数。
func CountOccurrences(s, sub string) int {
	return strings.Count(s, sub)
}

// Reverse 将字符串 s 反转。
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsIPv4 检查给定的字符串是否为有效的 IPv4 地址
func IsIPv4(ip string) bool {
	// 使用 net 包的 ParseIP 函数判断是否为合法的 IPv4 地址
	return net.ParseIP(ip) != nil && strings.Contains(ip, ".")
}

// IsIPv6 检查给定字符串是否为有效的 IPv6 地址
func IsIPv6(ip string) bool {
	// 使用 net 包的 ParseIP 函数判断是否为合法的 IPv6 地址
	return net.ParseIP(ip) != nil && strings.Count(ip, ":") >= 2
}

// IsURL 检查给定的字符串是否为有效的 URL
func IsURL(str string) (*url.URL, bool) {
	u, err := url.Parse(str)
	return u, err == nil && u.Scheme != "" && u.Host != ""
}

// IsIPOrURL 检查给定的字符串是否是有效的 IPv4IPv6 地址或 URL
func IsIPOrURL(str string) bool {
	_, b := IsURL(str)
	return IsIPv4(str) || IsIPv6(str) || b
}

// IsIPOrURLToBool 检查给定的字符串是否是有效的 IPv4IPv6 地址或 URL
func IsIPOrURLToBool(b1, b2, b3 bool) bool {
	return b1 || b2 || b3
}

// GetBasicAuth 设置权限Basic 参数
func GetBasicAuth(name, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(name+":"+password))
}

// SnakeToCamel 将下划线式的字符串转换为驼峰式的字符串
func SnakeToCamel(s string) string {
	// 创建一个空的字符串切片
	var result []string

	// 按照下划线分割字符串
	parts := strings.Split(s, "_")

	// 遍历每个部分
	for _, part := range parts {
		// 如果部分为空，跳过
		if part == "" {
			continue
		}
		// 将部分的首字母转换为大写，其余不变
		result = append(result, string(unicode.ToUpper(rune(part[0])))+part[1:])
	}

	// 将结果切片连接成一个字符串并返回
	return strings.Join(result, "")
}

// CamelToSnake 将驼峰式的字符串转换为下划线式的字符串
func CamelToSnake(s string) string {
	// 创建一个空的字符串切片
	var result []string

	// 遍历字符串中的每个字符
	for i, r := range s {
		// 如果是大写字母
		if unicode.IsUpper(r) {
			// 如果不是第一个字符，且前一个字符不是大写字母或下划线
			if i > 0 && !(unicode.IsUpper(rune(s[i-1])) || s[i-1] == '_') {
				// 在结果切片中添加一个下划线
				result = append(result, "_")
			}
			// 在结果切片中添加该字符的小写形式
			result = append(result, string(unicode.ToLower(r)))
		} else {
			// 如果不是大写字母，直接在结果切片中添加该字符
			result = append(result, string(r))
		}
	}

	// 将结果切片连接成一个字符串并返回
	return strings.Join(result, "")
}
