package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"bytes"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"endar_agent/config"
)

type Agent struct {
	ServerURL       string
	RegistrationKey string
	Interval        time.Duration
	DebugLevel      string
	ConfigPath      string
	PrintConfig     bool
	AID             string
	Registered      bool
	AgentInfo       AgentInfo
}

type AgentInfo struct {
	Key             string `json:"key"`
	Token           string `json:"token"`
	Version         string `json:"version"`
	Enabled         bool   `json:"enabled"`
	PublicAddr      string `json:"public_addr"`
	CPUCount        int    `json:"cpu_count"`
	LogicalCPUCount int    `json:"logical_cpu_count"`
	Memory          string `json:"memory"`
	Hostname        string `json:"hostname"`
	Family          string `json:"family"`
	Release         string `json:"release"`
	Uptime          int64  `json:"uptime"`
}

// Returns agent system info
func GetAgentInfo(registrationKey string) AgentInfo {
	hostname, _ := os.Hostname()
	hostInfo, _ := host.Info()
	memInfo, _ := mem.VirtualMemory()
	netInterfaces, _ := net.Interfaces()
	publicAddr := "Unknown"
	if len(netInterfaces) > 0 {
		publicAddr = netInterfaces[0].Addrs[0].Addr
	}

	return AgentInfo{
		Key:             hostname,
		Token:           registrationKey,
		Version:         "1.0.0",
		Enabled:         true,
		PublicAddr:      publicAddr,
		CPUCount:        runtime.NumCPU(),
		LogicalCPUCount: runtime.NumCPU(),
		Memory:          fmt.Sprintf("%.2f GB", float64(memInfo.Total)/1024/1024/1024),
		Hostname:        hostname,
		Family:          runtime.GOOS,
		Release:         hostInfo.Platform,
		Uptime:          int64(time.Since(time.Unix(int64(hostInfo.BootTime), 0)).Seconds()),
	}
}

// Registers the agent with the server
func (a *Agent) RegisterAgent() (error) {
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

	// If Agent is registered, save to config
	if resp.StatusCode == 200 || responseMap["registered"] == true {
		if agentID, exists := responseMap["agent_id"].(string); exists {
			a.AID = agentID
			a.Registered = true
			cfg, err := config.LoadConfig(a.ConfigPath)
			if err == nil {
				cfg.Section("configuration").Key("aid").SetValue(agentID)
				cfg.Section("configuration").Key("registered").SetValue("True")
				config.SaveConfig(a.ConfigPath, cfg)
				log.Println("Agent ID saved to config.")
			} else {
				log.Println("Error loading config to update AID:", err)
			}
		}
	}
	return nil
}
