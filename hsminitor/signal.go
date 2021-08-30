// Author stone-bird created on 2021/8/28 10:31.
// Email 1245863260@qq.com or g1245863260@gmail.com.
// Use of deal with singal
package hsminitor

import (
	"github.com/stonebirdjx/logger"
	"os"
	"os/signal"
	"syscall"
)

func HsSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	//signal.Notify(signalChan,
	//syscall.SIGHUP,
	//	syscall.SIGINT,
	//	syscall.SIGQUIT,
	//	syscall.SIGILL,
	//	syscall.SIGTRAP,
	//	syscall.SIGABRT,
	//	syscall.SIGBUS,
	//	syscall.SIGFPE,
	//	syscall.SIGKILL,
	//	syscall.SIGSEGV,
	//	syscall.SIGPIPE,
	//	syscall.SIGALRM,
	//	syscall.SIGTERM,
	//)
	for s := range signalChan {
		logger.GlobalLogger.Logger.Println(logger.Notice, "got signal:", s)
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL:
			logger.GlobalLogger.Logger.Println(logger.Warning, "got signal:", s, "program exit 1")
			os.Exit(1)
		}
	}
}