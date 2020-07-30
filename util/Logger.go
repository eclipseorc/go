package util

/******************************************************************************
Copyright:cloud
Author:cloudapex@126.com
Version:1.0
Date:2014-10-18
Description: 日志系统
******************************************************************************/
import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type logger struct {
	status ELoggerStatus

	level            ELogLevel
	outMode          ELogMode
	dirName          string
	fileName         string
	fileSuffix       string
	rotateMax        int
	rotateSize       int
	levelPrefixNames [ELL_Maxed]string
	filter           func(logLevel ELogLevel, msg string) bool

	fileSystmHandle *os.File
	fileSystmLogger *log.Logger
	fileLogicHandle *os.File
	fileLogiclogger *log.Logger

	chanMsgs chan lUnit
	chanExit chan int
	wgExit   sync.WaitGroup

	lastUpdateTime time.Time
}

func (this *logger) Init(conf *LogConf) *logger {
	this.status = ELS_Initing
	this.outMode = ELM_File

	this.level = UTD_LOG_LEVEL
	this.dirName = ExePathName()
	this.fileName = ExeName()
	this.fileSuffix = UTD_LOG_FILE_SUFFIX
	this.rotateMax = UTD_LOG_ROTATE_MAX
	this.rotateSize = UTD_LOG_ROTATE_SIZE

	this.levelPrefixNames = UTD_LOG_MSG_LV_PREFIXS
	this.chanMsgs = make(chan lUnit, UTD_LOG_CSIZE)
	this.chanExit = make(chan int)
	if conf != nil {
		Cast(conf.OutMode != 0, func() { this.outMode = conf.OutMode }, nil)
		this.level = ELogLevel(ShouldMax(int(conf.Level), int(ELL_Fatal), int(ELL_Error)))
		Cast(conf.DirName == "", func() { conf.DirName = this.dirName }, nil)
		Cast(strings.HasPrefix(conf.DirName, "./"), func() { this.dirName = path.Join(this.dirName, conf.DirName) }, func() { this.dirName = conf.DirName })
		Cast(conf.FileName != "", func() { this.fileName = conf.FileName }, nil)
		Cast(conf.FileSuffix != "", func() { this.fileSuffix = conf.FileSuffix }, nil)
		this.rotateMax = ShouldMin(conf.RotateMax, 0, UTD_LOG_ROTATE_MAX)
		this.rotateSize = ShouldMin(conf.RotateSize, 1*1024*1024-1, UTD_LOG_ROTATE_SIZE)
	}
	return this
}
func (this *logger) Start() {
	if this.status == ELS_Running {
		return
	}

	if BitHas(uint(this.outMode), uint(ELM_File)) {
		if err := os.MkdirAll(this.dirName, 0666); err != nil {
			panic(err)
		}

		path := fmt.Sprintf("%s/%s_%s.%s", this.dirName, this.fileName, "system", this.fileSuffix)

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		this.fileSystmHandle = file
		this.fileSystmLogger = log.New(file, "", log.LstdFlags)
		this.fileSystmLogger.Println("👌")

		this.fileLogicUpdate()
		this.fileLogiclogger.Println("·································START·································")
		this.fileLogiclogger.Println() // add space line
	}

	go this.loop()
	this.status = ELS_Running
}
func (this *logger) Stop() {
	if this.status != ELS_Running {
		return
	}
	this.chanExit <- 1
	time.Sleep(time.Millisecond * 10)
	this.status = ELS_Exiting
	close(this.chanMsgs)
	this.wgExit.Wait()
}

func (this *logger) GetLevel() ELogLevel { return this.level }

func (this *logger) SetLevel(lv ELogLevel) {
	this.level = ELogLevel(ShouldMax(int(lv), int(ELL_Fatal), int(ELL_Error)))
}

func (this *logger) UpdPrefix(lvPrefix [ELL_Maxed]string) {
	this.levelPrefixNames = lvPrefix
}
func (this *logger) UpdFilter(filter func(lv ELogLevel, msg string) bool) {
	this.filter = filter
}

func (this *logger) Debug(format string, v ...interface{}) {
	if this.canOutLog(ELL_Debug) {
		return
	}
	this.push(ELL_Debug, fmt.Sprintf(format+"\n", v...))
}
func (this *logger) Debugv(v ...interface{}) {
	if this.canOutLog(ELL_Debug) {
		return
	}
	this.push(ELL_Debug, fmt.Sprintln(v...))
}
func (this *logger) Trace(format string, v ...interface{}) {
	if this.canOutLog(ELL_Trace) {
		return
	}
	this.push(ELL_Trace, fmt.Sprintf(format+"\n", v...))
}
func (this *logger) Tracev(v ...interface{}) {
	if this.canOutLog(ELL_Trace) {
		return
	}
	this.push(ELL_Trace, fmt.Sprintln(v...))
}
func (this *logger) Info(format string, v ...interface{}) {
	if this.canOutLog(ELL_Infos) {
		return
	}
	this.push(ELL_Infos, fmt.Sprintf(format+"\n", v...))
}
func (this *logger) Infov(v ...interface{}) {
	if this.canOutLog(ELL_Infos) {
		return
	}
	this.push(ELL_Infos, fmt.Sprintln(v...))
}
func (this *logger) Warn(format string, v ...interface{}) {
	if this.canOutLog(ELL_Warns) {
		return
	}
	this.push(ELL_Warns, fmt.Sprintf(format+"\n", v...))
}
func (this *logger) Warnv(v ...interface{}) {
	if this.canOutLog(ELL_Warns) {
		return
	}
	this.push(ELL_Warns, fmt.Sprintln(v...))
}
func (this *logger) Error(format string, v ...interface{}) {
	if this.canOutLog(ELL_Error) {
		return
	}
	this.push(ELL_Error, fmt.Sprintf(format+"\n", v...)) // 输出到队列

	this.fileSystmLogger.Print(this.levelPrefixNames[ELL_Error] + " " + fmt.Sprintf(format+"\n", v...)) //直接输出到文件
}
func (this *logger) Errorv(v ...interface{}) {
	if this.canOutLog(ELL_Error) {
		return
	}
	this.push(ELL_Error, fmt.Sprintln(v...))

	this.fileSystmLogger.Print(this.levelPrefixNames[ELL_Error] + " " + fmt.Sprintln(v...)) //直接输出到文件
}
func (this *logger) Fatal(format string, v ...interface{}) {
	this.push(ELL_Fatal, fmt.Sprintf(format+"\n", v...))

	this.fileSystmLogger.Print(this.levelPrefixNames[ELL_Fatal] + " " + fmt.Sprintf(format+"\n", v...)) //直接输出到文件

	this.fileSystmHandle.Sync()
	os.Exit(1)
}
func (this *logger) Fatalv(v ...interface{}) {
	if len(v) == 0 || v[0] == nil {
		return
	}
	this.push(ELL_Fatal, fmt.Sprintln(v...))

	this.fileSystmLogger.Print(this.levelPrefixNames[ELL_Fatal] + " " + fmt.Sprintln(v...)) //直接输出到文件

	this.fileSystmHandle.Sync()
	os.Exit(1)
}

// --------  Internal logic
func (this *logger) push(level ELogLevel, msg string) {
	if this.status != ELS_Running {
		return
	}
	this.chanMsgs <- lUnit{level, this.levelPrefixNames[level] + " " + msg, time.Now()}
}
func (this *logger) canOutLog(lev ELogLevel) bool {
	if this.status != ELS_Running {
		return false
	}
	if lev < this.level {
		return true
	}
	return false
}
func (this *logger) fileLogicUpdate() {
	curTime := time.Now()
	if curTime.Year() == this.lastUpdateTime.Year() && curTime.Day() == this.lastUpdateTime.Day() {
		return
	}
	this.lastUpdateTime = curTime

	Cast(this.fileLogicHandle != nil, func() { this.fileLogicHandle.Close() }, nil)

	strTime := fmt.Sprintf("%04d-%02d-%02d", curTime.Year(), curTime.Month(), curTime.Day())

	path := fmt.Sprintf("%s/%s_%s.%s", this.dirName, this.fileName, strTime, this.fileSuffix)
	handle, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	this.fileLogicHandle = handle
	this.fileLogiclogger = log.New(handle, "", 0 /*log.Lmicroseconds*/)
}
func (this *logger) needRenameFile() bool {
	if this.rotateMax > 1 {
		if info, err := this.fileLogicHandle.Stat(); err == nil {
			return info.Size() >= int64(this.rotateSize)
		}
	}
	return false
}
func (this *logger) renameFile() {

	Cast(this.fileLogicHandle != nil, func() { this.fileLogicHandle.Close() }, nil)

	strTime := fmt.Sprintf("%04d-%02d-%02d", this.lastUpdateTime.Year(), this.lastUpdateTime.Month(), this.lastUpdateTime.Day())

	pathmax := fmt.Sprintf("%s/%s_%s.%d.%s", this.dirName, this.fileName, strTime, this.rotateMax, this.fileSuffix)

	Cast(FileExist(pathmax), func() { FileRemove(pathmax) }, nil)
	for index := this.rotateMax - 1; index > 0; index-- {
		pathOld := fmt.Sprintf("%s/%s_%s.%d.%s", this.dirName, this.fileName, strTime, index, this.fileSuffix)
		pathNew := fmt.Sprintf("%s/%s_%s.%d.%s", this.dirName, this.fileName, strTime, index+1, this.fileSuffix)
		Cast(FileExist(pathOld), func() { FileRename(pathOld, pathNew) }, nil)
	}

	pathNow := fmt.Sprintf("%s/%s_%s.%s", this.dirName, this.fileName, strTime, this.fileSuffix)
	pathOne := fmt.Sprintf("%s/%s_%s.%d.%s", this.dirName, this.fileName, strTime, 1, this.fileSuffix)
	FileRename(pathNow, pathOne)

	file, err := os.OpenFile(pathNow, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	this.fileLogicHandle = file
	this.fileLogiclogger.SetOutput(file)
}
func (this *logger) loop() {
	this.wgExit.Add(1)
	count, t := 0, time.NewTicker(time.Second*1)

	defer func() {
		Cast(this.fileSystmHandle != nil, func() {
			this.fileSystmLogger.Println("✋")
			this.fileSystmHandle.Close()
		}, nil)

		Cast(this.fileLogicHandle != nil, func() {
			this.fileLogiclogger.Println("··································END··································")
			this.fileLogicHandle.Close()
		}, nil)

		this.status = ELS_Stopped
		if x := recover(); x != nil {
			Cast(this.fileSystmHandle != nil, func() { this.fileSystmLogger.Print("caught panic in logger::loop() error:", x) }, nil)
		}
		this.wgExit.Done()
	}()

	for {
		select {
		case msg := <-this.chanMsgs:
			if BitHas(uint(this.outMode), uint(ELM_Std)) || msg.l >= ELL_Infos {
				Print(msg.s)
			}
			if this.filter != nil && this.filter(msg.l, msg.s) {
				continue
			}
			if !BitHas(uint(this.outMode), uint(ELM_File)) {
				continue
			}
			this.fileLogicUpdate()
			Cast(count%UTD_LOG_DTM_ONCE == 0, func() {
				count = 1
				this.fileLogiclogger.Printf("%s [DTM] %v\n\n", time.Now().Format("15:04:05.000"), time.Now().Format("2006-01-02 15:04:05.000"))
			}, func() { count++ })
			this.fileLogiclogger.Println(msg.t.Format("15:04:05.000") + " " + msg.s)
		case <-t.C:
			if this.fileLogicHandle != nil {
				this.fileLogicHandle.Sync()
				Cast(this.needRenameFile(), func() { this.renameFile() }, nil)
			}
			if this.fileSystmHandle != nil {
				this.fileSystmHandle.Sync()
			}
		case <-this.chanExit:
			for msg := range this.chanMsgs {
				if BitHas(uint(this.outMode), uint(ELM_File)) {
					this.fileLogiclogger.Println(msg.t.Format("15:04:05.000") + " " + msg.s)
				}
			}
			return
		}
	}
}
