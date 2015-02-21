package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	configuration "github.com/Bo0mer/os-agent/config"
	"github.com/Bo0mer/os-agent/facade"
	"github.com/Bo0mer/os-agent/jobstore"
	. "github.com/Bo0mer/os-agent/masterclient"
	"github.com/Bo0mer/os-agent/model"
	"github.com/Bo0mer/os-agent/server"

	l4g "code.google.com/p/log4go"
	"github.com/Bo0mer/executor"
)

func main() {
	config := loadConfiguration()
	s := createServer(config.Server)

	osAgentFacade := facade.NewOSAgentFacade(executor.NewExecutor(), jobstore.NewJobStore())

	createJobHandler := server.NewHandler("POST", "/jobs", osAgentFacade.CreateJob)
	getJobHandler := server.NewHandler("GET", "/jobs", osAgentFacade.GetJob)

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
		Id:   config.Id,
		Host: config.Host,
		Port: config.Port,
	}

	c := NewMasterClient(config.Master.URL, self)
	stop := make(chan struct{})
	go sendHeartbeat(config.Master.HeartbeatInterval, c, stop)

	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	<-osChan
	stop <- struct{}{}

	l4g.Info("Shtutting down HTTP server...")
	s.Stop()
	l4g.Info("Shutdown successsful.")
}

func loadConfiguration() configuration.OSAgentConfig {
	configDir := os.Getenv("OS_AGENT_CONFIG_DIR")
	configFile := fmt.Sprintf("%s%s", configDir, "/config.yml")
	config, err := configuration.LoadConfig(configFile)
	if err != nil {
		l4g.Error("Could not load configuration. Error: %s", err)
		panic("Could not load configuration!")
	}
	return config
}

func createServer(conf configuration.ServerConfig) server.Server {
	s := server.NewServer(conf.Host, conf.Port)
	if conf.Auth.Enabled {
		simpleAuthenticator := createSimpleAuthenticator(conf.Auth.User, conf.Auth.Password)
		s.SetAuthenticator(simpleAuthenticator)
	}
	return s
}

func createSimpleAuthenticator(username, password string) server.Authenticator {
	authFunc := func(un, pwd string) bool {
		if un == username && pwd == password {
			return true
		}
		return false
	}

	return server.NewSimpleAuthenticator(authFunc)
}

func sendHeartbeat(interval int, c MasterClient, stop <-chan struct{}) {
	sendRegister(c)
	for {
		select {
		case <-time.After(time.Second * time.Duration(interval)):
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
