package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

const (
	_VER string = "1.0.0"
)

// LEVEL define logger level
type LEVEL int32

var logLevel LEVEL = 1
var maxFileSize int64
var maxFileCount int32

var dailyRolling = true
var consoleAppender = true

// RollingFile ...
var RollingFile = false
var logObj *FILE

// DATEFORMAT is logger file format
const DATEFORMAT = "2006-01-02"

// UNIT be defined int64 type
type UNIT int64

// every logger file size
const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

// Logger level
const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

// FILE is the logger file struct
type FILE struct {
	dir      string
	filename string
	_suffix  int
	isCover  bool
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

func init() {
	SetRollingDaily("log", "application.log")
}

// SetConsole use to control logger output
// false: can not output the logger to terminate,
// true: display the logger into the terminate
func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

// SetLevel set logger level
func SetLevel(_level LEVEL) {
	logLevel = _level
}

// SetRollingFile split the logger file by fileSize
func SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) error {
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	RollingFile = true
	dailyRolling = false
	logObj = &FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()
	for i := 1; i <= int(maxNumber); i++ {
		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
			logObj._suffix = i
		} else {
			break
		}
	}
	var err error
	if !logObj.isMustRename() {
		logObj.logfile, err = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		logObj.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}
	go fileMonitor()
	return nil
}

// SetRollingDaily splits the logger file by date format
func SetRollingDaily(fileDir, fileName string) error {
	RollingFile = false
	dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	logObj = &FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()
	var err error
	if !logObj.isMustRename() {
		logObj.logfile, err = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		logObj.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}
	return nil
}

// Debug Level Logger
func Debug(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()

	if logLevel <= DEBUG {
		if logObj.lg != nil {
			logObj.lg.Output(2, formatOutStr("[debug] ", v))
		}
		console(formatOutStr("[debug] ", v))
	}
}

func Debugf(format string, v ... interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()

	if logLevel <= DEBUG {
		if logObj.lg != nil {
			logObj.lg.Output(2, fmt.Sprintf("[debug] " + format, v ...))
		}
		console(fmt.Sprintf("[debug] " + format, v ...))
	}
}

// Info Level Logger
func Info(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	// fmt.Printf("%#v\r\n", logObj)
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()

	if logLevel <= INFO {
		if logObj.lg != nil {
			logObj.lg.Output(2, formatOutStr("[info] ", v))
		}
		console(formatOutStr("[info] ", v))
	}
}

// Info Level Logger
func Infof(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	// fmt.Printf("%#v\r\n", logObj)
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()

	if logLevel <= INFO {
		if logObj.lg != nil {
			logObj.lg.Output(2, fmt.Sprintf("[info] " + format, v ...))
		}
		console(fmt.Sprintf("[info] " + format, v ...))
	}
}


// Warn Level Logger
func Warn(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	if logLevel <= WARN {
		if logObj.lg != nil {
			logObj.lg.Output(2, formatOutStr("[warn] ", v))
		}
		console(formatOutStr("[warn] ", v))
	}
}

// Warn Level Logger
func Warnf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	if logLevel <= WARN {
		if logObj.lg != nil {
			logObj.lg.Output(2, fmt.Sprintf("[warn] " + format, v ...))
		}
		console(fmt.Sprintf("[warn] " + format, v ...))
	}
}

// Error Level Logger
func Error(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	if logLevel <= ERROR {
		if logObj.lg != nil {
			logObj.lg.Output(2, formatOutStr("[error] ", v))
		}
		console(formatOutStr("[error] ", v))
	}
}

// Error Level Logger
func Errorf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	if logLevel <= ERROR {
		if logObj.lg != nil {
			logObj.lg.Output(2, fmt.Sprintf("[error] " + format, v ...))
		}
		console(fmt.Sprintf("[error] " + format, v ...))
	}
}

// Fatal Level Logger
// the level will cause application exit
func Fatal(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	defer func() {
		os.Exit(-127)
	}()
	if logLevel <= FATAL {
		if logObj.lg != nil {
			logObj.lg.Output(2, formatOutStr("[fatal] ", v))
		}
		console(formatOutStr("[fatal] ", v))
	}
}

// Fatal Level Logger
// the level will cause application exit
func Fatalf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	logObj.mu.RLock()
	defer logObj.mu.RUnlock()
	defer func() {
		os.Exit(-127)
	}()
	if logLevel <= FATAL {
		if logObj.lg != nil {
			logObj.lg.Output(2, fmt.Sprintf("[fatal] " + format, v ...))
		}
		console(fmt.Sprintf("[fatal] " + format, v ...))
	}
}


func (f *FILE) isMustRename() bool {
	if dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	} else {
		if maxFileCount > 1 {
			if fileSize(f.dir+"/"+f.filename) >= maxFileSize {
				return true
			}
		}
	}
	return false
}

func (f *FILE) rename() {
	if dailyRolling {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
			f.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
		}
	} else {
		f.coverNextOne()
	}
}

func (f *FILE) nextSuffix() int {
	return int(f._suffix%int(maxFileCount) + 1)
}

func (f *FILE) coverNextOne() {
	f._suffix = f.nextSuffix()
	if f.logfile != nil {
		f.logfile.Close()
	}
	if isExist(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix))) {
		os.Remove(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix)))
	}
	os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+"."+strconv.Itoa(int(f._suffix)))
	f.logfile, _ = os.Create(f.dir + "/" + f.filename)
	f.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			fileCheck()
		}
	}
}

func formatOutStr(logMsg ...interface{}) string {
	logStr := logMsg[0].(string)
	userMsg := logMsg[1].([]interface{})
	for _, v := range userMsg {
		logStr = fmt.Sprintf("%v%v ", logStr, v)
	}
	return logStr
}

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
	}
}

func console(s string) {
	if consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(file+":"+strconv.Itoa(line), s)
	}
}

func catchError() {
	if err := recover(); err != nil {
		debug.PrintStack()
	}
}
