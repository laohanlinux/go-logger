package example

import (
	"os"
	"runtime"
	"strconv"
	"testing"

	logger "github.com/laohanlinux/go-logger/logger"
)

func log(i int) {
	//	logger.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Info(strconv.Itoa(i))
	logger.Warn(strconv.Itoa(i))
	logger.Error(strconv.Itoa(i))
	logger.Debug(strconv.Itoa(i))
	logger.Debugf("%s", strconv.Itoa(i))
	logger.Debug(nil)
	logger.Info("canot ")
	logger.Infof("%s", "你好啊")
	//logger.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
}

func TestNil(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.SetConsole(true)
	for i := 0; i < 1024; i++ {
		logger.Info(i)
		logger.Info(i, i+1)
	}
}

func TestAll(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//指定是否控制台打印，默认为true
	logger.SetConsole(false)
	//指定日志文件备份方式为文件大小的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//第三个参数为备份文件最大数量
	//第四个参数为备份文件大小
	//第五个参数为文件大小的单位
	//logger.SetRollingFile("d:/logtest", "test.log", 10, 5, logger.KB)

	//指定日志文件备份方式为日期的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	logSavePath := "../example"
	stat, err := os.Stat(logSavePath)
	if err != nil {
		panic(err)
	}
	if !stat.IsDir() {
		panic("example is not a dir")
	}
	logger.SetRollingDaily(logSavePath, "test.log")

	//指定日志级别  ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF 级别由低到高
	//一般习惯是测试阶段为debug，生成环境为info以上
	logger.SetLevel(logger.DEBUG)

	for i := 100; i > 0; i-- {
		log(i)
	}
	//time.Sleep(15 * time.Second)
}
