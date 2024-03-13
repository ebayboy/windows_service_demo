package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	fasttick := time.Tick(500 * time.Millisecond)
	slowtick := time.Tick(2 * time.Second)
	tick := fasttick
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case <-tick:
			// 在这里执行定期的服务操作
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				debug.Logf("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func main() {
	if len(svcFlag) != 0 {
		runService(serviceName, false)
		return
	}
	runService(serviceName, true)
}

func runService(name string, isDebug bool) {
	var elog debug.Log
	if isDebug {
		elog = debug.New(name)
	} else {
		var err error
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, "starting...")
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err := run(name, &myService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("service failed: %v", err))
		return
	}
	elog.Info(1, "stopped")
}
