package performance

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"endar_agent/agent"
)

// Struct for performance data matching the server model
type AgentPerformance struct {
	CPULoad        float64            `json:"cpu_load"`
	CPULoadPerCore map[string]float64 `json:"cpu_load_per_core"`
	MemTotal       float64            `json:"mem_total"`
	MemUsed        float64            `json:"mem_used"`
	MemFree        float64            `json:"mem_free"`
	MemPercentUsed float64            `json:"mem_percent_used"`
	MemTotalH      string             `json:"mem_total_h"`
	MemUsedH       string             `json:"mem_used_h"`
	MemFreeH       string             `json:"mem_free_h"`
	SwapTotal      float64            `json:"swap_total"`
	SwapUsed       float64            `json:"swap_used"`
	SwapFree       float64            `json:"swap_free"`
	SwapPercent    float64            `json:"swap_percent_used"`
	SwapTotalH     string             `json:"swap_total_h"`
	SwapUsedH      string             `json:"swap_used_h"`
	SwapFreeH      string             `json:"swap_free_h"`
	SwappedIn      uint64            `json:"swapped_in"`
	SwappedOut     uint64            `json:"swapped_out"`
}

// Converts bytes to human-readable format (GB)
func formatBytes(bytes uint64) string {
	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
}

// Collects system performance metrics
func CollectPerformanceData() (AgentPerformance, error) {
	// CPU usage
	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return AgentPerformance{}, err
	}

	// Memory usage
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return AgentPerformance{}, err
	}

	// Swap memory
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return AgentPerformance{}, err
	}

	// Convert CPU Load Per Core into a map
	cpuCoreLoadMap := make(map[string]float64)
	for i, load := range cpuPercent {
		cpuCoreLoadMap[fmt.Sprintf("core_%d", i)] = load
	}

	// Return the structured performance data
	return AgentPerformance{
		CPULoad:        cpuPercent[0],
		CPULoadPerCore: cpuCoreLoadMap,
		MemTotal:       float64(memInfo.Total),
		MemUsed:        float64(memInfo.Used),
		MemFree:        float64(memInfo.Free),
		MemPercentUsed: memInfo.UsedPercent,
		MemTotalH:      formatBytes(memInfo.Total),
		MemUsedH:       formatBytes(memInfo.Used),
		MemFreeH:       formatBytes(memInfo.Free),
		SwapTotal:      float64(swapInfo.Total),
		SwapUsed:       float64(swapInfo.Used),
		SwapFree:       float64(swapInfo.Free),
		SwapPercent:    swapInfo.UsedPercent,
		SwapTotalH:     formatBytes(swapInfo.Total),
		SwapUsedH:      formatBytes(swapInfo.Used),
		SwapFreeH:      formatBytes(swapInfo.Free),
		SwappedIn:      swapInfo.Sin,
		SwappedOut:     swapInfo.Sout,
	}, nil
}

// Compress data using zlib
func compressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

// Upload performance data
func PostPerformanceData(a *agent.Agent, data AgentPerformance) error {
	// Wrap data inside "name" and "data" keys
	payload, err := json.Marshal(map[string]interface{}{
		"name": "get-performance",
		"data": data,
	})
	if err != nil {
		return err
	}

	// Compress the payload
	compressedPayload, err := compressData(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agent/collection", a.ServerURL), bytes.NewBuffer(compressedPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Encoding", "zlib")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("tenant-key", a.RegistrationKey)
	req.Header.Set("aid", a.AID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to upload performance data, status code: %d\n", resp.StatusCode)
		return fmt.Errorf("failed to upload performance data, status code: %d", resp.StatusCode)
	}

	return nil
}
