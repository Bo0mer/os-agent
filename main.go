package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Bo0mer/os-agent/facade"
	"github.com/Bo0mer/os-agent/jobstore"
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
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	<-osChan
	l4g.Info("Shtutting down HTTP server...")
	s.Stop()
	l4g.Info("Shutdown successsful.")
}
