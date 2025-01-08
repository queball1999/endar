package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"gopkg.in/ini.v1"
)

type SystemInfo struct {
	Key             string `json:"key"`
	Token           string `json:"token"`
	Version         string `json:"version"`
	Enabled         bool   `json:"enabled"`
	InstallGroup    string `json:"install_group"`
	PublicAddr      string `json:"public_addr"`
	Uninstall       bool   `json:"uninstall"`
	CPUCount        int    `json:"cpu_count"`
	LogicalCPUCount int    `json:"logical_cpu_count"`
	Memory          string `json:"memory"`
	Hostname        string `json:"hostname"`
	FQDN            string `json:"fqdn"`
	Family          string `json:"family"`
	Release         string `json:"release"`
	SysVersion      string `json:"sys_version"`
	InstallType     string `json:"install_type"`
	Edition         string `json:"edition"`
	Build           string `json:"build"`
	Machine         string `json:"machine"`
	Processor       string `json:"processor"`
	LastBoot        string `json:"last_boot"`
	Uptime          int64  `json:"uptime"`
}

type Agent struct {
	ServerURL       string
	RegistrationKey string
	Interval        time.Duration
	DebugLevel      string
	ConfigPath      string
	PrintConfig     bool
	AID             string
	Registered      bool
	AgentInfo       SystemInfo
}

type program struct {
	agent *Agent
}

func loadConfig(configPath string) (*ini.File, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config.ini: %v", err)
	}
	return cfg, nil
}

func saveConfig(configPath string, cfg *ini.File) error {
	return cfg.SaveTo(configPath)
}

func saveDefaultConfig(configPath, serverURL, registrationToken string) error {
	cfg := ini.Empty()
	cfg.Section("configuration").Key("server").SetValue(serverURL)
	cfg.Section("configuration").Key("key").SetValue(registrationToken)
	cfg.Section("configuration").Key("interval").SetValue("300")
	cfg.Section("configuration").Key("level").SetValue("DEBUG")
	cfg.Section("configuration").Key("print_config").SetValue("False")
	cfg.Section("configuration").Key("aid").SetValue("")
	cfg.Section("configuration").Key("registered").SetValue("False")
	return saveConfig(configPath, cfg)
}

func initializeConfig() (string, string, int, string, string, bool, string, bool) {
	serverURL := flag.String("server", "", "Server URL")
	registrationToken := flag.String("key", "", "Registration Token")
	configPath := flag.String("config", "config.ini", "Path to configuration file")
	flag.Parse()

	//FIXME: Double check the flags are overwritten when config loaded

	var cfg *ini.File
	var err error

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Println("No config file found, creating new one...")
		if err := saveDefaultConfig(*configPath, *serverURL, *registrationToken); err != nil {
			log.Fatalf("Failed to save default config: %v", err)
		}
	}

	cfg, err = loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server := cfg.Section("configuration").Key("server").String()
	key := cfg.Section("configuration").Key("key").String()
	interval, _ := cfg.Section("configuration").Key("interval").Int()
	debugLevel := cfg.Section("configuration").Key("level").String()
	configFilePath := cfg.Section("configuration").Key("config_path").String()
	printConfig, _ := cfg.Section("configuration").Key("print_config").Bool()
	aid := cfg.Section("configuration").Key("aid").String()
	registered, _ := cfg.Section("configuration").Key("registered").Bool()

	if *serverURL != "" {
		server = *serverURL
		cfg.Section("configuration").Key("server").SetValue(*serverURL)
	}
	if *registrationToken != "" {
		key = *registrationToken
		cfg.Section("configuration").Key("key").SetValue(*registrationToken)
	}
	if *configPath != "" {
		configFilePath = *configPath
		cfg.Section("configuration").Key("config_path").SetValue(*configPath)
	}

	saveConfig(*configPath, cfg)

	if printConfig {
		fmt.Println("Server URL:", server)
		fmt.Println("Registration Token:", key)
		fmt.Println("Interval:", interval)
		fmt.Println("Debug Level:", debugLevel)
		fmt.Println("Config Path:", configFilePath)
		fmt.Println("Print Config:", printConfig)
		fmt.Println("AID:", aid)
		fmt.Println("Registered:", registered)
	}

	return server, key, interval, debugLevel, configFilePath, printConfig, aid, registered
}


func getLogicalCPUCount() int {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0
	}
	return count
}

func getSystemInfo(registrationKey string) SystemInfo {
	hostname, _ := os.Hostname()
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	netInterfaces, _ := net.Interfaces()
	publicAddr := "Unknown"
	if len(netInterfaces) > 0 {
		publicAddr = netInterfaces[0].Addrs[0].Addr
	}

	lastBoot := time.Unix(int64(hostInfo.BootTime), 0).Format("2006-01-02 15:04:05")

	return SystemInfo{
		Key:             hostname,
		Token:           registrationKey,
		Version:         "1.0.0",
		Enabled:         true,
		InstallGroup:    "Default Group",
		PublicAddr:      publicAddr,
		Uninstall:       false,
		CPUCount:        len(cpuInfo),
		LogicalCPUCount: getLogicalCPUCount(),
		Memory:          fmt.Sprintf("%.2f GB", float64(memInfo.Total)/1024/1024/1024),
		Hostname:        hostname,
		FQDN:            hostname,
		Family:          runtime.GOOS,
		Release:         hostInfo.Platform,
		SysVersion:      hostInfo.PlatformVersion,
		InstallType:     "Client",
		Edition:         "Unknown",
		Build:           hostInfo.KernelVersion,
		Machine:         hostInfo.PlatformFamily,
		Processor:       cpuInfo[0].ModelName,
		LastBoot:        lastBoot,
		Uptime:          int64(time.Since(time.Unix(int64(hostInfo.BootTime), 0)).Seconds()),
	}
}

func (a *Agent) registerAgent() (error) {
	log.Println("Registering agent...")
	data, _ := json.Marshal(a.AgentInfo)
	log.Println("Data:", string(data))
	
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agent/register", a.ServerURL), bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("tenant-key", a.RegistrationKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error connecting to server:", err)
		return err
	}
	defer resp.Body.Close()

	log.Println("Status Code:", resp.StatusCode)
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		log.Println("Error decoding response:", err)
		return err
	}
	log.Println("Response From Server:", responseMap)

	// Exit if the tenant token is invalid
	if resp.StatusCode == 424 && responseMap["msg"] == "invalid tenant token" {
		log.Println("Invalid tenant token, exiting...")
		os.Exit(1)
		
	}

	// If successfully registered, save to config
	if resp.StatusCode == 200 {
		if agentID, exists := responseMap["agent_id"].(string); exists {
			a.AID = agentID
			cfg, err := loadConfig(a.ConfigPath)
			if err == nil {
				cfg.Section("configuration").Key("aid").SetValue(agentID)
				cfg.Section("configuration").Key("registered").SetValue("True")
				saveConfig(a.ConfigPath, cfg)
				log.Println("Agent ID saved to config.")
			} else {
				log.Println("Error loading config to update AID:", err)
			}
		}
	}
	return nil
}

func (p *program) Start(s service.Service) error {
	log.Println("Service is starting...")
	go p.run()
	return nil
}

func (p *program) run() {
	// This is the main loop of the agent
	// Place logic here to run at the specified interval
	for {
		fmt.Println("AID:", p.agent.AID, "Registered:", p.agent.Registered)
		if p.agent.AID == "" || !p.agent.Registered {
			error := p.agent.registerAgent()	// Attempt to register the agent with the server
			if error != nil {
				log.Println("Error registering agent:", error)
				os.Exit(1)
			}
		}

		log.Println("Sleeping for", p.agent.Interval)
		time.Sleep(p.agent.Interval)
	}
}

func (p *program) Stop(s service.Service) error {
	log.Println("Service is stopping...")
	return nil
}

func main() {
	// Load configuration settings
	server, key, interval, debugLevel, configPath, printConfig, aid, registered := initializeConfig()

	agent := &Agent{
		ServerURL:       server,
		RegistrationKey: key,
		Interval:        time.Duration(interval) * time.Second,
		DebugLevel:      debugLevel,
		ConfigPath:		 configPath,
		PrintConfig:     printConfig,
		AID:             aid,
		Registered:      registered,
		AgentInfo:       getSystemInfo(key),
	}

	svcConfig := &service.Config{
		Name:        "EndarAgent",
		DisplayName: "Endar Monitoring Agent",
		Description: "A cross-platform system monitoring agent.",
	}

	prg := &program{agent: agent}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Command-line flag parsing
	// FIXME: Ensure these work correctly
    if len(os.Args) > 1 {
        command := os.Args[1]
		log.Println("Command:", command)
        switch command {
		case "install":
			err := s.Install()
			if err != nil {
				log.Fatal("Failed to install service:", err)
			}
			fmt.Println("Service installed successfully.")
			return
		case "uninstall":
			err := s.Uninstall()
			if err != nil {
				log.Fatal("Failed to uninstall service:", err)
			}
			fmt.Println("Service uninstalled successfully.")
			return
		case "start":
			err := s.Start()
			if err != nil {
				log.Fatal("Failed to start service:", err)
			}
			fmt.Println("Service started successfully.")
			return
		case "stop":
			err := s.Stop()
			if err != nil {
				log.Fatal("Failed to stop service:", err)
			}
			fmt.Println("Service stopped successfully.")
			return
		case "restart":
			err := s.Restart()
			if err != nil {
				log.Fatal("Failed to restart service:", err)
			}
			fmt.Println("Service restarted successfully.")
			return
		case "status":
            status, err := s.Status()
            if err != nil {
                fmt.Println("Could not determine service status. It may not be installed.")
                return
            }
            switch status {
            case service.StatusRunning:
                fmt.Println("Service is running.")
            case service.StatusStopped:
                fmt.Println("Service is stopped.")
            default:
                fmt.Println("Service status unknown.")
            }
            return
        default:
            fmt.Println("Unknown command:", command)
            fmt.Println("Usage: ./endar-server [install|uninstall|start|stop|restart|status]")
            return
		}
	}

    err = s.Run()
    if err != nil {
        log.Fatal(err)
    }
}