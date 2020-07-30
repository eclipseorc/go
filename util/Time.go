package util

import (
	"time"
)

/*
 * 描述: 字符转换为时间戳
 *
 * 	strTime : 格式为 "2018-01-10"
 *
 ***************************************************************************/
func StringToTime(strTime string) int64 {
	//获取本地location
	strTem := "2006-01-02"
	timeLocal, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(strTem, strTime, timeLocal)
	return theTime.Unix()
}

/*
 * 描述: 整数转换为时间戳
 *
 ***************************************************************************/
func IntToTime(nTimer int64) time.Time {
	return time.Unix(nTimer, 0)
}

/*
 * 描述: 字符转换为时间戳
 *
 * 	strTime : 格式为 "2018-01-10 03:04:05"
 *
 ***************************************************************************/
func StringToTimeEx(strTime string) int64 {
	//获取本地location
	strTem := "2006-01-02 15:04:05"
	timeLocal, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(strTem, strTime, timeLocal)
	return theTime.Unix()
}

/*
 * 描述: 时间戳转格式化为字符串
 *
 ***************************************************************************/
func TimeToString(nTimer int64) string {
	tm := time.Unix(nTimer, 0)
	//strTime := tm.Format("2006-01-02 03:04:05 PM")
	return tm.Format("2006-01-02")
}

/*
 * 描述: 时间戳转格式化为字符串
 *
 ***************************************************************************/
func TimeToStringEx(nTimer int64) string {
	tm := time.Unix(nTimer, 0)
	return tm.Format("2006-01-02 03:04:05 PM")
}

/*
 * 描述: 查看 输入时间戳 是否是今天
 *
 * 	checkTime : 查看的时间戳
 *
 ***************************************************************************/
func IsToday(checkTime int64) bool {

	nowTime := time.Now().Unix()
	nowTime -= (nowTime + 28800) % 86400

	if checkTime > nowTime {
		return true
	}
	return false
}

/*
 * 描述: 获取输入 时间戳 的零点时间戳
 *
 * 	checkTime : 所求的时间戳
 *
 ***************************************************************************/
func GetZero(checkTime int64) int64 {
	checkTime -= (checkTime + 28800) % 86400
	return checkTime
}

/*
 * 描述: 整点，向下取整
 *
 * 	checkTime : 所求的时间戳
 *
 ***************************************************************************/
func GetHourZero(checkTime uint32) uint32 {
	return checkTime - (checkTime % 3600)
}

/*
 * 描述: 求所当前月1号到下个月1号，两个时间戳
 *
 ***************************************************************************/
func Interval() (int64, int64) {
	year := time.Now().Year()
	month := time.Now().Month()
	stime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	etime := time.Date(year, (month + 1), 1, 0, 0, 0, 0, time.Local).Unix()
	return stime, etime
}
