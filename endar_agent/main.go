package main

import (
	"log"
	"time"
	"os"
	"io"
	"path/filepath"

	"github.com/kardianos/service"
	"endar_agent/agent"
	"endar_agent/config"
	"endar_agent/performance"
	"endar_agent/disk"
)

type program struct {
	agent *agent.Agent
}

// Initialize logging to a file in the same directory as the script
func setupLogging() {
	// Get the directory of the running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	execDir := filepath.Dir(execPath)
	logFilePath := filepath.Join(execDir, "endar_agent.log")
	log.Println("Log file path:", logFilePath)

	// Open log file for writing (append mode)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("Logging initialized. Writing to:", logFilePath)
}

// Start service
func (p *program) Start(s service.Service) error {
	log.Println("Service is starting...")
	go p.run()
	return nil
}

// Main loop of the agent
func (p *program) run() {
	for {
		log.Println("Agent is running. AID:", p.agent.AID, "Registered:", p.agent.Registered)
		if p.agent.AID == "" && !p.agent.Registered {
			err := p.agent.RegisterAgent()
			if err != nil {
				log.Println("Error registering agent:", err)
				return
			}
		} 
		
		if p.agent.AID != "" && p.agent.Registered {
			// Collect and post performance data
			performanceData, err := performance.CollectPerformanceData()
			if err != nil {
				log.Println("Error collecting performance data:", err)
			} else {
				log.Println("Performance data:", performanceData)
				err = performance.PostPerformanceData(p.agent, performanceData)
				if err != nil {
					log.Println("Error posting performance data:", err)
				} else {
					log.Println("Performance data posted successfully.")
				}
			}

			// Collect and post disk data
			diskData, err := disk.CollectDiskData()
			if err != nil {
				log.Println("Error collecting disk data:", err)
			} else {
				//log.Println("Disk data:", diskData)
				err = disk.PostDiskData(p.agent, diskData)
				if err != nil {
					log.Println("Error posting disk data:", err)
				} else {
					log.Println("Disk data posted successfully.")
				}
			}
		}

		log.Println("Sleeping for", p.agent.Interval)
		time.Sleep(p.agent.Interval)
	}
}

// Stop service
func (p *program) Stop(s service.Service) error {
	log.Println("Service is stopping...")
	return nil
}

// Handle service commands (install, start, stop, restart, status)
func handleServiceCommands(s service.Service, flags map[string]interface{}) {
	var err error

	switch {
	case flags["install"].(bool):
		err = s.Install()
		if err != nil {
			log.Fatal("Failed to install service:", err)
		}
		log.Println("Service installed successfully.\nRun 'endar_agent.exe --start' to start the service.")
		os.Exit(0)

	case flags["uninstall"].(bool):
		log.Println("Stopping service...")
		err = s.Stop()
		if err != nil {
			log.Println("Service may not be running or installed, proceeding with uninstallation.")
		}

		err = s.Uninstall()
		if err != nil {
			log.Fatal("Failed to uninstall service:", err)
		}
		log.Println("Service uninstalled successfully.")
		os.Exit(0)

	case flags["start"].(bool):
		err = s.Start()
		if err != nil {
			log.Fatal("Failed to start service:", err)
		}
		log.Println("Service started successfully.\nRun 'endar_agent.exe --stop' to stop the service.")
		os.Exit(0)

	case flags["stop"].(bool):
		err = s.Stop()
		if err != nil {
			log.Fatal("Failed to stop service:", err)
		}
		log.Println("Service stopped successfully.\nRun 'endar_agent.exe --start' to start the service.")
		os.Exit(0)

	case flags["restart"].(bool):
		err = s.Restart()
		if err != nil {
			log.Fatal("Failed to restart service:", err)
		}
		log.Println("Service restarted successfully.\nRun 'endar_agent.exe --status' to check the service status.")
		os.Exit(0)

	case flags["status"].(bool):
		status, err := s.Status()
		if err != nil {
			log.Println("Could not determine service status. It may not be installed.\nRun 'endar_agent.exe --install' to install the service.")
			os.Exit(0)
		}
		switch status {
		case service.StatusRunning:
			log.Println("Service is running.")
		case service.StatusStopped:
			log.Println("Service is stopped.")
		default:
			log.Println("Service status unknown.")
		}
		os.Exit(0)
	}
}

func main() {
	// Setup logging before anything else
	setupLogging()

	// Parse command-line flags
	flags := config.ParseFlags()

	// Load configuration settings
	server, key, interval, debugLevel, configPath, printConfig, aid, registered := config.InitializeConfig(flags)

	// Create agent instance
	ag := &agent.Agent{
		ServerURL:       server,
		RegistrationKey: key,
		Interval:        time.Duration(interval) * time.Second,
		DebugLevel:      debugLevel,
		ConfigPath:      configPath,
		PrintConfig:     printConfig,
		AID:             aid,
		Registered:      registered,
		AgentInfo:       agent.GetAgentInfo(key),
	}

	// Define service configuration
	svcConfig := &service.Config{
		Name:        "endar_agent",
		DisplayName: "Endar Monitoring Agent",
		Description: "A cross-platform system monitoring agent.",
	}

	// Create the service program
	prg := &program{agent: ag}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Handle service commands
	handleServiceCommands(s, flags)

	// Start the service
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
