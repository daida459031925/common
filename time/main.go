package time

import "time"

/**
时间结构体
*/
type Date struct {
	t time.Time
}

func (date *Date) Init() (Date, error) {
	d := Date{
		t: time.Now(),
	}
	return d, nil
}

/**
获取日期
*/
func (date *Date) getDay() int {
	return date.t.Day()
}

/**
获取分钟
*/
func (date *Date) getMinute() int {
	return time.Now().Minute()
}

/**
获取月份
*/
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

/**
获取年分
*/
func GetYear() {
	time.Now().Year()
}
