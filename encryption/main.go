package encryption

import (
	"fmt"
	"regexp"
	"strconv"
)

// 密码强度等级，D为最低
const (
	LevelD level = 0
	LevelC level = 1
	LevelB level = 2
	LevelA level = 3
	LevelS level = 4
)

type level int

func (l level) getValue() int {
	i, e := strconv.Atoi(fmt.Sprintf("%s", l))
	if e != nil {
		return 0
	}
	return i
}

/*
 *  minLength: 指定密码的最小长度
 *  maxLength：指定密码的最大长度
 *  minLevel：指定密码最低要求的强度等级
 *  pwd：用户输入的明文密码
 */
func Check(minLength, maxLength int, minLevel level, inPwd string) error {
	// 首先校验密码长度是否在范围内
	if len(inPwd) < minLength {
		return fmt.Errorf("BAD PASSWORD: The password is shorter than %d characters", minLength)
	}
	if len(inPwd) > maxLength {
		return fmt.Errorf("BAD PASSWORD: The password is logner than %d characters", maxLength)
	}

	// 初始化密码强度等级为D，利用正则校验密码强度，若匹配成功则强度自增1
	var level int = LevelD.getValue()
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, inPwd)
		if match {
			level++
		}
	}

	// 如果最终密码强度低于要求的最低强度，返回并报错
	if level < minLevel.getValue() {
		return fmt.Errorf("The password does not satisfy the current policy requirements. ")
	}
	return nil
}
