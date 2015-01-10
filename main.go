package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bo0mer/os-agent/facade"
	"github.com/Bo0mer/os-agent/jobstore"
	"github.com/Bo0mer/os-agent/server"

	"github.com/Bo0mer/executor"
)

func main() {

	osAgentFacade := facade.NewOSAgentFacade(executor.NewExecutor(), jobstore.NewJobStore())

	createJobHandler := server.NewHandler("POST", "/jobs", osAgentFacade.CreateJob)
	getJobHandler := server.NewHandler("GET", "/jobs", osAgentFacade.GetJob)

	s := server.NewServer("127.0.0.1", 8080)
	s.Register(createJobHandler)
	s.Register(getJobHandler)

	fmt.Println("Starting...")
	err := s.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Start successful.")
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	<-osChan
	fmt.Println("Shtutting down...")
	s.Stop()
	fmt.Println("Shutdown successsful.")
}
