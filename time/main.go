package time

import "time"

// Date 时间结构体
type Date struct {
	t time.Time
}

func (date *Date) Init() (Date, error) {
	d := Date{
		t: time.Now(),
	}
	return d, nil
}

// 获取日期
func (date *Date) getDay() int {
	return date.t.Day()
}

// 获取分钟
func (date *Date) getMinute() int {
	return time.Now().Minute()
}

// 获取月份
func (date *Date) getMonth() time.Month {
	return time.Now().Month()
}

func monthToInt(month string) {
	switch month {
	case "January":
		break
	case "February":
		break
	case "March":
		break
	case "April":
		break
	case "May":
		break
	case "June":
		break
	case "July":
		break
	case "August":
		break
	case "September":
		break
	case "October":
		break
	case "November":
		break
	case "December":
		break
	default:
		break
	}

}

// GetYear 获取年分
func GetYear() {
	time.Now().Year()
}

// GetSecond 传入秒获取对应time.Duration类型
func GetSecond(second int) time.Duration {
	return time.Duration(second) * time.Second
}

// GetMinute 传入分钟获取对应time.Duration类型
func GetMinute(minute int) time.Duration {
	return time.Duration(minute) * time.Minute
}

// GetHour 传入小时获取对应time.Duration类型
func GetHour(hour int) time.Duration {
	return time.Duration(hour) * time.Hour
}

// SetSleepSecond 传入秒对应线程休眠
func SetSleepSecond(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

// SetSleepMinute 传入分钟对应线程休眠
func SetSleepMinute(minute int) {
	time.Sleep(time.Duration(minute) * time.Minute)
}

// SetSleepHour 传入小时获对应线程休眠
func SetSleepHour(hour int) {
	time.Sleep(time.Duration(hour) * time.Hour)
}
