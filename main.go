package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bo0mer/os-agent/facade"
	"github.com/Bo0mer/os-agent/jobstore"
	. "github.com/Bo0mer/os-agent/masterclient"
	"github.com/Bo0mer/os-agent/model"
	"github.com/Bo0mer/os-agent/server"

	l4g "code.google.com/p/log4go"
	"github.com/Bo0mer/executor"
)

func main() {

	osAgentFacade := facade.NewOSAgentFacade(executor.NewExecutor(), jobstore.NewJobStore())

	createJobHandler := server.NewHandler("POST", "/jobs", osAgentFacade.CreateJob)
	getJobHandler := server.NewHandler("GET", "/jobs", osAgentFacade.GetJob)

	s := server.NewServer("127.0.0.1", 8080)
	s.Register(createJobHandler)
	s.Register(getJobHandler)

	l4g.Info("Starting HTTP server...")
	err := s.Start()
	if err != nil {
		l4g.Error("Unable to start server", err)
		return
	}
	l4g.Info("Start successful.")

	self := model.Slave{
		Id:   "unique-id",
		Host: "127.0.0.1",
		Port: 8080,
	}

	c := NewMasterClient("http://127.0.0.1", self)
	stop := make(chan struct{})
	go sendHeartbeat(c, stop)

	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	<-osChan
	stop <- struct{}{}

	l4g.Info("Shtutting down HTTP server...")
	s.Stop()
	l4g.Info("Shutdown successsful.")
}

func sendHeartbeat(c MasterClient, stop <-chan struct{}) {
	for {
		select {
		case <-time.After(time.Minute * 5):
			sendRegister(c)
		case <-stop:
			l4g.Debug("Stop heartbeat sending.")
			return
		}
	}
}

func sendRegister(c MasterClient) {
	l4g.Debug("Sending heartbeat to master...")
	err := c.Register()
	if err != nil {
		l4g.Error("Error registering on master: %s", err.Error())
		return
	}
	l4g.Info("Heartbeat successfuly sent.")
}
