package log

import (
	"github.com/astaxie/beego/logs"
	"sync"
)

var (
	gBgLog        *logs.BeeLogger
	onceInitBgLog sync.Once
)

func initLog() {
	gBgLog = logs.NewLogger(10000) // 创建一个日志记录器，参数为缓冲区的大小
}

//get log instance
func Instance() *logs.BeeLogger {
	onceInitBgLog.Do(func() {
		initLog()
	})
	return gBgLog
}
