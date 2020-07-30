package util

/******************************************************************************
Copyright:cloud
Author:cloudapex@126.com
Version:1.0
Date:2014-10-18
Description:通用函数
******************************************************************************/
import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	muRand sync.Mutex
	seRand *rand.Rand

	timeFormats []string

	muGuid                  sync.Mutex
	sequence, lastTimestamp int64
	lastGId                 TGUID
)

func init() {
	seRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	timeFormats = []string{
		"2006-01-02", "2006-01-02 15", "2006-01-02 15:04", "2006-01-02 15:04:05",
		"1/2/2006", "1/2/2006 15", "1/2/2006 15:4", "1/2/2006 15:4:5",
		"01-02", "02 15:04:05", "Mon 15:04:05", "15:04:05", "15:04", "15", // 每年/每月/每周/每日(秒,分,时)
		"15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST"}
}

// -------- Common
func Args(params ...interface{}) []interface{} { return params }

func Cast(condition bool, trueFun, falseFun func()) {
	if condition {
		if trueFun != nil {
			trueFun()
		}
	} else {
		if falseFun != nil {
			falseFun()
		}
	}
}
func Call(condition bool, trueFun, falseFun func() interface{}) interface{} {
	if condition {
		if trueFun != nil {
			return trueFun()
		}
	} else {
		if falseFun != nil {
			return falseFun()
		}
	}
	return nil
}
func Clone(des, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(des)
}

// -------- random
func RandInt(max int) int {
	if max <= 0 {
		return 0
	}
	defer UnLock(Lock(&muRand))
	value := seRand.Int()
	return value % max
}
func RandFnt() float32 {
	defer UnLock(Lock(&muRand))
	return seRand.Float32()
}
func RandGID() TGUID {
	defer UnLock(Lock(&muGuid))

	var workerID int64 = 1
	ts := time.Now().UnixNano() >> 20

	if ts < lastTimestamp {
		return 0
	}

	if lastTimestamp == ts {
		sequence = (sequence + 1) & UTD_RANDOM_SEQUENCE_MASK
		if sequence == 0 {
			return 0
		}
	} else {
		sequence = 0
	}

	lastTimestamp = ts

	id := TGUID(((ts - UTD_RANDOM_TWEPOCH) << UTD_RANDOM_TIMESTAMP_SHIFT) |
		(workerID << UTD_RANDOM_WORKERID_SHIFT) | sequence)

	if id <= lastGId {
		return 0
	}
	lastGId = id
	return id
}
func RandStr(num int) string { // 1-0 and a - Z
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, num)

	for i := range b {
		b[i] = letterRunes[RandInt(len(letterRunes))]
	}
	return string(b)
}
func RandInts(max, num int) []int { // duplicate
	var l []int
	for n := 0; n < num; n++ {
		l = append(l, RandInt(max))
	}
	return l
}
func RandIntDis(max, num int) []int { // no duplicate
	l, m, num := []int{}, map[int]int{}, Min(max, num)
	for n := 0; n < num; {
		v := RandInt(max)
		if _, ok := m[v]; ok {
			continue
		}
		l = append(l, v)
		n++
	}
	return l
}

// -------- base type limit
func Sum(arr ...int) int {
	num := 0
	for _, val := range arr {
		num += val
	}
	return num
}
func Sumf(arr ...float32) float32 {
	num := float32(0)
	for _, val := range arr {
		num += val
	}
	return num
}
func Min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
func Minf(f1, f2 float32) float32 {
	if f1 < f2 {
		return f1
	}
	return f2
}
func Max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}
func Maxf(f1, f2 float32) float32 {
	if f1 > f2 {
		return f1
	}
	return f2
}
func MinMax(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
func MinMaxf(value, minf, maxf float32) float32 {
	if value < minf {
		return minf
	}
	if value > maxf {
		return maxf
	}
	return value
}
func ShouldMin(value, minVal, should int) int {
	if value <= minVal {
		return should
	}
	return value
}
func ShouldMax(value, maxVal, should int) int {
	if value >= maxVal {
		return should
	}
	return value
}
func ShouldStr(value, desStr, should string) string {
	if value == desStr {
		return should
	}
	return value
}

// -------- carry and Patch bit
func Carry(lv, val, incrVal, valRef int32) (new_lv, new_val int32) {
	val += incrVal
	if incrVal >= 0 {
		if val >= valRef {
			lv++
			val -= valRef
		}
	} else {
		if val <= 0 {
			lv--
			val += valRef
		}
	}
	return lv, val
}

// -------- check
func Intv(vars []int) int {
	if len(vars) > 0 {
		return vars[0]
	}
	return 0
}
func Boolv(vars []bool) bool {
	return len(vars) > 0 && vars[0]
}
func Between(value, min, max int) bool {
	if value < min || value > max {
		return false
	}
	return true
}
func Contain(val int, vals ...int) bool {
	for _, it := range vals {
		if val == it {
			return true
		}
	}
	return false
}
func Containi(val int32, vals ...int32) bool {
	for _, it := range vals {
		if val == it {
			return true
		}
	}
	return false
}
func ContainI(val int64, vals ...int64) bool {
	for _, it := range vals {
		if val == it {
			return true
		}
	}
	return false
}
func Contains(val string, vals ...string) bool {
	for _, it := range vals {
		if val == it {
			return true
		}
	}
	return false
}

// ContainS 不区分大小写
func ContainS(val string, vals ...string) bool {
	for _, it := range vals {
		if strings.EqualFold(val, it) {
			return true
		}
	}
	return false
}

// -------- slice

// Insert 给任意切片类型插入元素
func Insert(slice interface{}, index int, value interface{}) (interface{}, bool) {
	// 判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	// 参数检查
	if index < 0 || index > v.Len() || reflect.TypeOf(slice).Elem() != reflect.TypeOf(value) {
		return nil, false
	}

	// 尾部追加元素
	if index == v.Len() {
		return reflect.Append(v, reflect.ValueOf(value)).Interface(), true
	}

	// 插入位置赋值
	v = reflect.AppendSlice(v.Slice(0, index+1), v.Slice(index, v.Len()))
	v.Index(index).Set(reflect.ValueOf(value))
	return v.Interface(), true
}

// Delete 删除任意切片类型指定下标的元素
func Delete(slice interface{}, index int) (interface{}, bool) {
	// 判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}
	// 参数检查
	if v.Len() == 0 || index < 0 || index > v.Len()-1 {
		return nil, false
	}
	// 删除元素
	return reflect.AppendSlice(v.Slice(0, index), v.Slice(index+1, v.Len())).Interface(), true
}

// Update 修改任意切片类型指定下标的元素
func Update(slice interface{}, index int, value interface{}) (interface{}, bool) {
	//判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	//参数检查
	if index > v.Len()-1 || reflect.TypeOf(slice).Elem() != reflect.TypeOf(value) {
		return nil, false
	}

	// 更新位置赋值
	v.Index(index).Set(reflect.ValueOf(value))

	return v.Interface(), true
}

// Search 查找指定元素在任意切片类型中的所有下标
func search(slice interface{}, value interface{}) ([]int, bool) {
	//判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	// 查找元素
	var idxs []int
	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == reflect.ValueOf(value).Interface() {
			idxs = append(idxs, i)
		}
	}
	return idxs, true
}

// -------- bits
func BitHas(value uint, flags ...uint) bool {
	for _, flag := range flags {
		if value&flag == 0 {
			return false
		}
	}
	return true
}
func BitSet(value uint, flags ...uint) uint {
	for _, flag := range flags {
		value |= flag
	}
	return value
}
func BitDel(value uint, flags ...uint) uint {
	for _, flag := range flags {
		value &^= flag
	}
	return value
}

// -------- time duration
func DurationMic(mic int) time.Duration {
	return time.Duration(mic) * time.Microsecond
}
func DurationMil(mil int) time.Duration {
	return time.Duration(mil) * time.Millisecond
}
func DurationSec(sec int) time.Duration {
	return time.Duration(sec) * time.Second
}

// -------- time datatime
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
func TimeParsed(datum time.Time, datetime string) (time.Time, error) {
	var t time.Time
	Cast(datum.IsZero(), func() { datum = time.Now() }, nil)

	parseTime := []int{}
	datumTime := []int{datum.Second(), datum.Minute(), datum.Hour(), datum.Day(), int(datum.Month()), datum.Year()}
	datumLocation := datum.Location()

	var err error
	onlyTime := regexp.MustCompile(`^\s*\d+(:\d+)*\s*$`).MatchString(datetime) // match 15:04:05, 15

	var idx = -1
	var formatOk bool
	for n, format := range timeFormats {
		if t, err = time.Parse(format, datetime); err == nil {
			idx = n
			formatOk = true
			break
		}
	}
	if !formatOk {
		return t, fmt.Errorf("Can't parse string as time: %s", datetime)
	}

	location := t.Location()
	if location.String() == "UTC" {
		location = datumLocation
	}

	parseTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
	//onlyTime = onlyTime && (parseTime[3] == 1) && (parseTime[4] == 1)

	for i := len(parseTime) - 1; i >= 0; i-- {
		if onlyTime && i == 3 {
			break
		}
		if i == 4 && strings.Contains(timeFormats[idx], "01") {
			break
		}
		if i == 3 && strings.Contains(timeFormats[idx], "02") {
			break
		}
		if parseTime[i] > 0 && !(i == 3 || i == 4) {
			break
		}
		parseTime[i] = datumTime[i]
	}

	if len(parseTime) > 0 {
		t = time.Date(parseTime[5], time.Month(parseTime[4]), parseTime[3], parseTime[2], parseTime[1], parseTime[0], 0, location)
		// = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
		if strings.Contains(timeFormats[idx], "Mon") {
			week := WeekDayEn2n(datetime[:3])
			if week == -1 {
				return t, fmt.Errorf("Can't parse string as time: %s", datetime)
			}
			if diff := week - t.Weekday(); diff >= 0 {
				t = t.Add(24 * time.Hour * time.Duration(diff))
			} else {
				t = t.Add(24 * time.Hour * time.Duration(diff+7))
			}
		}
	}
	return t, err
}
func TimeParseBase(datum time.Time, datetime string) time.Time {
	t, _ := TimeParsed(datum, datetime)
	return t
}
func TimeParseFull(datetime string) time.Time {
	return TimeParseBase(time.Time{}, datetime)
}
func WeekDayEn2n(shortWeekEn string) time.Weekday {
	switch shortWeekEn {
	case "Sun":
		return time.Sunday
	case "Mon":
		return time.Monday
	case "Tue":
		return time.Tuesday
	case "Wed":
		return time.Wednesday
	case "Thu":
		return time.Thursday
	case "Fri":
		return time.Friday
	case "Sat":
		return time.Saturday
	}
	return -1
}

// -------- mutex lock
func AutoLock(mu *sync.Mutex, fun func()) {
	mu.Lock()
	fun()
	mu.Unlock()
}
func AutoRLock(mu *sync.RWMutex, fun func()) {
	mu.RLock()
	fun()
	mu.RUnlock()
}
func AutoWLock(mu *sync.RWMutex, fun func()) {
	mu.Lock()
	fun()
	mu.Unlock()
}
func Lock(mu *sync.Mutex) *sync.Mutex { mu.Lock(); return mu }

func UnLock(lockFun *sync.Mutex) { lockFun.Unlock() }

func RLock(mu *sync.RWMutex) *sync.RWMutex { mu.RLock(); return mu }

func RUnLock(rMutex *sync.RWMutex) { rMutex.RUnlock() }

func WLock(mu *sync.RWMutex) *sync.RWMutex { mu.Lock(); return mu }

func WUnLock(wMutex *sync.RWMutex) { wMutex.Unlock() }
