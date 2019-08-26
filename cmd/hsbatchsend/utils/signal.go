package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/orientwalt/htdf/cmd/hsbatchsend/log"
)

func WaitSignal(signalHandler func(os.Signal)) {
	log.Instance().Info("process runing...")
	csignal := make(chan os.Signal, 1)
	signal.Notify(csignal, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)
	for {
		s := <-csignal

		signalHandler(s)
	}
}
