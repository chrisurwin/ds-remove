package agent

import (
	"chrisurwin/alerting-agent/healthcheck"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

type Agent struct {
	sync.WaitGroup

	probePeriod time.Duration
}

func NewAgent(probePeriod time.Duration) *Agent {
	return &Agent{
		probePeriod: probePeriod,
		//log:         log.WithField("pkg", "agent"),
	}
}

func (a *Agent) Start() {
	go healthcheck.StartHealthcheck()

	t := time.NewTicker(a.probePeriod)
	for _ = range t.C {
		go a.checkDiskSpace()

		a.Wait()
	}
}

func (a *Agent) checkDiskSpace() {
	a.Add(1)
	var stat syscall.Statfs_t

	wd, err := os.Getwd()

	syscall.Statfs(wd, &stat)

	// Available blocks * size per block = available space in bytes
	fmt.Println(stat.Bavail * uint64(stat.Bsize))
}
