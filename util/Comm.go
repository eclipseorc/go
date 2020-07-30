package util

/******************************************************************************
Copyright:cloud
Author:cloudapex@126.com
Version:1.0
Date:2014-10-18
Description: util库相关定义
******************************************************************************/
import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

/*------------------------------------------------------------------------------
=======================================常量定义==================================
------------------------------------------------------------------------------*/
const (
	UTD_VERSION                = "1.0.0"          // ulern库版本号
	UTD_LOG_LEVEL              = ELL_Debug        // 默认日志等级
	UTD_LOG_FILE_SUFFIX        = "log"            // 默认日志文件后缀名
	UTD_LOG_ROTATE_MAX         = 3                // 默认日志文件轮换数量
	UTD_LOG_ROTATE_SIZE        = 20 * 1024 * 1024 // 默认日志文件轮换size
	UTD_LOG_CSIZE              = 100              // 默认日志消息通道缓存大小
	UTD_LOG_DTM_ONCE           = 50               // 逻辑日志每写多少条,写一次[DTM]日期记录
	UTD_RANDOM_WORKERID_BITS   = uint64(10)
	UTD_RANDOM_SEQUENCE_BITS   = uint64(12)
	UTD_RANDOM_WORKERID_SHIFT  = UTD_RANDOM_SEQUENCE_BITS
	UTD_RANDOM_TIMESTAMP_SHIFT = UTD_RANDOM_SEQUENCE_BITS + UTD_RANDOM_WORKERID_BITS
	UTD_RANDOM_SEQUENCE_MASK   = int64(-1) ^ (int64(-1) << UTD_RANDOM_SEQUENCE_BITS)
	UTD_RANDOM_TWEPOCH         = int64(1288834974288) // ( 2012-10-28 16:23:42 UTC ).UnixNano() >> 20
	UTD_RANDOM_CSIZE           = 100
)

var (
	UTD_LOG_MSG_LV_PREFIXS = [ELL_Maxed]string{"[TRC]", "[DBG]", "[INF]", "[WRN]", "[ERR]", "[FAIL]"}
)

/*------------------------------------------------------------------------------
=====================================接口定义====================================
------------------------------------------------------------------------------*/

// ILoger interface
type ILoger interface {
	// TRACE
	Trace(format string, v ...interface{})
	Tracev(v ...interface{})

	// DEBUG
	Debug(format string, v ...interface{})
	Debugv(v ...interface{})

	// INFO
	Info(format string, v ...interface{})
	Infov(v ...interface{})

	// WARN
	Warn(format string, v ...interface{})
	Warnv(v ...interface{})

	// ERROR
	Error(format string, v ...interface{})
	Errorv(v ...interface{})
}

// ILogMe interface
type ILogMe interface {
	ILoger

	WarnEnil(format string, v ...interface{}) error
	ErrorEnil(format string, v ...interface{}) error

	// FATAL
	Fatal(format string, v ...interface{})
	Fatalv(v ...interface{})
}

/*------------------------------------------------------------------------------
=====================================类型定义====================================
------------------------------------------------------------------------------*/

// 时间段(别名)
type TTime = time.Time
type TDurt = time.Duration

// 时间点(别名)
type TJsTime time.Time // Support json.Unmarshaler and json.Marshaler and fmt.Stringer
func (t *TJsTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", TimeFormat(time.Time(*t)))), nil
}
func (t *TJsTime) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("TJsTime UnmarshalJSON len(data) < 2")
	}
	the, err := TimeParsed(time.Now(), string(data[1:len(data)-1])) // 引号去掉
	*t = TJsTime(the)
	return err
}
func (t TJsTime) String() string {
	return TimeFormat(time.Time(t))
}

// 日志等级
type ELogLevel int //
const (
	ELL_Trace ELogLevel = iota
	ELL_Debug
	ELL_Infos
	ELL_Warns
	ELL_Error
	ELL_Fatal
	ELL_Maxed // 6
) //
func (e ELogLevel) String() string {
	if e >= ELL_Trace && e < ELL_Maxed {
		return UTD_LOG_MSG_LV_PREFIXS[e]
	}
	return fmt.Sprintf("ELL_Unkonw(%d)", e)
}

// 日志运行状态
type ELoggerStatus int //
const (
	ELS_Initing ELoggerStatus = iota
	ELS_Running
	ELS_Exiting
	ELS_Stopped
	ELS_Max
) //
func (e ELoggerStatus) String() string {
	switch e {
	case ELS_Initing:
		return "Initing"
	case ELS_Running:
		return "Running"
	case ELS_Exiting:
		return "Exiting"
	case ELS_Stopped:
		return "Stopped"
	}
	return fmt.Sprintf("ELS_Unkonw(%d)", e)
}

// 日志输出模式
type ELogMode int //
const (
	ELM_Std ELogMode = 1 << iota
	ELM_File
	ELM_Max
) //
func (e ELogMode) String() string {
	var str = []string{}
	if BitHas(uint(e), uint(ELM_Std)) {
		str = append(str, "Std")
	}
	if BitHas(uint(e), uint(ELM_File)) {
		str = append(str, "File")
	}
	return strings.Join(str, "+")
}

// GUID TYPE
type TGUID uint64 //
func (t TGUID) Hex() []byte {
	var h [16]byte
	var b [8]byte

	b[0] = byte(t >> 56)
	b[1] = byte(t >> 48)
	b[2] = byte(t >> 40)
	b[3] = byte(t >> 32)
	b[4] = byte(t >> 24)
	b[5] = byte(t >> 16)
	b[6] = byte(t >> 8)
	b[7] = byte(t)

	len := hex.Encode(h[:], b[:])
	return h[:len]
}
func (t TGUID) HexStr() string { return string(t.Hex()) }

/*------------------------------------------------------------------------------
=====================================结构定义====================================
------------------------------------------------------------------------------*/

// 日志单元
type lUnit struct {
	l ELogLevel
	s string
	t time.Time
}

// 日志配置
type LogConf struct {
	Level      ELogLevel // 日志等级[ELL_Debug]
	OutMode    ELogMode  // 日志输出模式
	DirName    string    // 输出目录[默认在程序所在目录]
	FileName   string    // 日志文件主名[程序本身名]
	FileSuffix string    // 日志文件后缀[log]
	RotateMax  int       // 日志文件轮换数量[3]
	RotateSize int       // 日志文件轮换大小[20m]
}

// 日志自己(带堆栈)
type LogMe struct {
	theme string
	name  func() string
}

func (l *LogMe) Init(theme string, name func() string) { l.theme, l.name = theme, name }

func (l *LogMe) Trace(format string, v ...interface{}) {
	if !CanOutLog(ELL_Trace) {
		return
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Trace("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Tracev(v ...interface{}) {
	if !CanOutLog(ELL_Trace) {
		return
	}
	msg := fmt.Sprint(v...)
	file, line, fun := Stack()
	Log.Trace("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Debug(format string, v ...interface{}) {
	if !CanOutLog(ELL_Debug) {
		return
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Debug("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Debugv(v ...interface{}) {
	if !CanOutLog(ELL_Debug) {
		return
	}
	msg := fmt.Sprint(v...)
	file, line, fun := Stack()
	Log.Debug("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Info(format string, v ...interface{}) {
	if !CanOutLog(ELL_Infos) {
		return
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Info("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Infov(v ...interface{}) {
	if !CanOutLog(ELL_Infos) {
		return
	}
	msg := fmt.Sprint(v...)
	file, line, fun := Stack()
	Log.Info("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Warn(format string, v ...interface{}) {
	if !CanOutLog(ELL_Warns) {
		return
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Warn("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Warnv(v ...interface{}) {
	if !CanOutLog(ELL_Warns) {
		return
	}
	msg := fmt.Sprint(v...)
	file, line, fun := Stack()
	Log.Warn("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) WarnEnil(format string, v ...interface{}) error {
	if !CanOutLog(ELL_Warns) {
		return nil
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Warn("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
	return nil
}
func (l *LogMe) Error(format string, v ...interface{}) {
	if !CanOutLog(ELL_Error) {
		return
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Error("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Errorv(v ...interface{}) {
	if !CanOutLog(ELL_Error) {
		return
	}
	msg := fmt.Sprint(v...)
	file, line, fun := Stack()
	Log.Error("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) ErrorEnil(format string, v ...interface{}) error {
	if !CanOutLog(ELL_Error) {
		return nil
	}
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Error("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
	return nil
}
func (l *LogMe) Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Fatal("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
func (l *LogMe) Fatalv(v ...interface{}) {
	msg := fmt.Sprint(v...)
	file, line, fun := Stack()
	Log.Fatal("%s[%q]%s:%d|%s() %s", l.theme, l.name(), file, line, fun, msg)
}
