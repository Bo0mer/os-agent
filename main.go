package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Bo0mer/os-agent/facade"
	"github.com/Bo0mer/os-agent/server"

	"github.com/Bo0mer/executor"
)

func main() {

	osAgentFacade := facade.NewOSAgentFacade(executor.NewExecutor())

	commandHandler := server.NewHandler("POST", "/command", osAgentFacade.ExecuteCommand)

	s := server.NewServer("127.0.0.1", 8080)
	s.Register(commandHandler)

	s.Start()

	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	<-osChan
}
