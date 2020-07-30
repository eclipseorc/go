package util

/******************************************************************************
Copyright:cloud
Author:cloudapex@126.com
Version:1.0
Date:2014-10-18
Description:系统函数
******************************************************************************/
import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Log *logger // Main Log

	chanSig     = make(chan os.Signal)
	chanExit    = make(chan int)
	startup     = time.Now()
	exitHandles = []func(){}
	screen      = log.New(os.Stdout, "", log.LstdFlags)
	ppWebSrv    *http.Server
)

func init() {
	Sign(func() {
		Cast(Log != nil, func() {
			Log.Stop()
		}, nil)
	})
}
func Init(lconf ...*LogConf) {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	signal.Notify(chanSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	go func() {
		<-chanSig
		for i := len(exitHandles) - 1; i >= 0; i-- {
			exitHandles[i]()
		}

		exitCode := 0
		select {
		case chanExit <- exitCode:
			return
		case <-time.After(time.Second * 1):
			os.Exit(exitCode)
		}
	}()

	switch {
	case len(lconf) == 0: // nothing
		Log = (&logger{}).Init(&LogConf{OutMode: ELM_Std})
	case len(lconf) >= 0 && lconf[0] == nil: // nil point
		Log = (&logger{}).Init(nil)
	default:
		Log = (&logger{}).Init(lconf[0])
	}
	Log.Start()
}
func Exec(do func()) {
	Init()
	do()
	Quit()
	Wait()
}
func Wait(x ...interface{}) {
	<-chanExit
}
func Sign(hand func()) {
	exitHandles = append(exitHandles, hand)
}
func Quit(delay ...time.Duration) {
	go func() {
		if len(delay) > 0 && delay[0] > 0 {
			<-time.After(delay[0])
		}
		close(chanSig)
	}()
}
func Catch(desc string, x interface{}, bFatal ...bool) bool {
	if x == nil {
		return false
	}
	head := fmt.Sprintf("%s. recover:%v\n", desc, x)

	buf := make([]byte, 256*7)
	size := runtime.Stack(buf, true)
	stack := string(buf[0:size])
	Cast(len(bFatal) > 0 && bFatal[0], func() { Log.Fatal("%s %s", head, stack) }, func() { Log.Error("%s %s", head, stack) })
	return true
}
func GoId() (int, error) {
	var buf [64]byte
	idField := strings.Fields(strings.TrimPrefix(string(buf[:runtime.Stack(buf[:], false)]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return 0, fmt.Errorf("cannot get goroutine id: %v", err)
	}
	return id, nil
}
func Stack() (file string, line int, fun string) {
	for index := 2; index < 5; index++ {
		_pc, _file, _line, _ := runtime.Caller(index)
		file = filepath.Base(_file)
		f := runtime.FuncForPC(_pc)
		fields := strings.Split(f.Name(), ".")
		line, fun = _line, fields[len(fields)-1]
		if _, err := strconv.ParseInt(strings.TrimPrefix(fun, "func"), 10, 32); err == nil {
			continue
		}
		if fun == "Cast" || fun == "Call" {
			continue
		}
		break
	}
	return
}
func StackStr() string {
	file, line, fun := "", 0, ""
	for index := 2; index < 5; index++ {
		_pc, _file, _line, _ := runtime.Caller(index)
		file = filepath.Base(_file)
		f := runtime.FuncForPC(_pc)
		fields := strings.Split(f.Name(), ".")
		line, fun = _line, fields[len(fields)-1]
		if _, err := strconv.ParseInt(strings.TrimPrefix(fun, "func"), 10, 32); err == nil {
			continue
		}
		if fun == "Cast" || fun == "Call" {
			continue
		}
		break
	}
	return fmt.Sprintf("%s:%d|%s()", file, line, fun)
}
func TimeStart() time.Time {
	return startup
}
func TimeLived() time.Duration {
	return time.Since(startup)
}
func LookCmd(cmdName string) bool {
	if _, err := exec.LookPath(cmdName); err != nil {
		return false
	}
	return true
}
func ExecCmd(cmd string, wait bool, arg ...string) (string, error) {
	if !LookCmd(cmd) {
		return "", fmt.Errorf("ExecCommand not found:%s", cmd)
	}
	c := exec.Command(cmd, arg...)
	if !wait {
		if err := c.Start(); err != nil {
			return "", err
		}
		return "", nil
	}
	out := bytes.NewBuffer(nil)
	c.Stderr = out
	if err := c.Run(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}
func Ping(ip string) error {
	cmd := exec.Command("ping", "-n", "1", ip)
	out := bytes.NewBuffer(nil)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return err
	}
	if strings.Contains(out.String(), "=32") {
		return nil
	}
	return fmt.Errorf("Target host unavailable")
}
func StartPprof(addr string) {
	if ppWebSrv != nil {
		return
	}
	ppWebSrv = &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}
	go func() {
		if err := ppWebSrv.ListenAndServe(); err != nil {
			Log.Warn("Pprof Server err:%v", err)
			ppWebSrv = nil
		}
	}()
}
func StopPprof() {
	if ppWebSrv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		ppWebSrv.Shutdown(ctx)
		cancel()
		ppWebSrv = nil
	}
}
func Goroutine(name string, goFun func()) {
	go func() {
		defer func() { Catch(fmt.Sprintf("Goroutine[%s] crash", name), recover()) }()
		goFun()
	}()
}
func ExeName() string {
	return strings.Split(ExeFullName(), ".")[0]
}
func ExeFullName() string {
	return filepath.Base(os.Args[0])
}
func ExeFullPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}
func ExePathName() string {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return path + "/"
}
func ExePathJoin(file string) string {
	return path.Join(ExePathName(), file)
}
func FunSelfName(fun interface{}) string {
	return FunFullName(fun, '.')
}
func FunFullName(fun interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()

	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return strings.Split(fields[size-1], "-")[0]
	}
	return ""
}
func CanOutLog(lv ELogLevel) bool {
	return Log.GetLevel() <= lv
}

func Print(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	screen.Print(s)
}

func Trace(format string, v ...interface{}) {
	if !CanOutLog(ELL_Trace) {
		return
	}
	str := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Trace("%s:%d|%s() %s", file, line, fun, str)
}
func Debug(format string, v ...interface{}) {
	if !CanOutLog(ELL_Debug) {
		return
	}
	str := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Debug("%s:%d|%s() %s", file, line, fun, str)
}
func Info(format string, v ...interface{}) {
	if !CanOutLog(ELL_Infos) {
		return
	}
	str := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Info("%s:%d|%s() %s", file, line, fun, str)
}
func Warn(format string, v ...interface{}) {
	if !CanOutLog(ELL_Warns) {
		return
	}
	str := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Warn("%s:%d|%s() %s", file, line, fun, str)
}
func Error(format string, v ...interface{}) {
	if !CanOutLog(ELL_Error) {
		return
	}
	str := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Error("%s:%d|%s() %s", file, line, fun, str)
}
func Fatal(format string, v ...interface{}) {
	if !CanOutLog(ELL_Fatal) {
		return
	}
	str := fmt.Sprintf(format, v...)
	file, line, fun := Stack()
	Log.Fatal("%s:%d|%s() %s", file, line, fun, str)
}
