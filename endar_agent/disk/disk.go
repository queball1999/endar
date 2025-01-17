package disk

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/disk"
	"endar_agent/agent"
)

type AgentDisk struct {
	UsedPercent int    `json:"used_percent"`
	Used        string `json:"used"`
	Mount       string `json:"mount"`
	Free        string `json:"free"`
	FSType      string `json:"fs_type"`
	Device      string `json:"device"`
	Total       string `json:"total"`
	Options     string `json:"options"`
}

// Converts bytes to human-readable format (GB)
func formatBytes(bytes uint64) string {
	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
}

// Collects disk usage information
func CollectDiskData() ([]AgentDisk, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	var disks []AgentDisk
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // Skip this partition if there's an error
		}

		disks = append(disks, AgentDisk{
			UsedPercent: int(usage.UsedPercent),
			Used:        formatBytes(usage.Used),
			Mount:       p.Mountpoint,
			Free:        formatBytes(usage.Free),
			FSType:      p.Fstype,
			Device:      p.Device,
			Total:       formatBytes(usage.Total),
			Options:     fmt.Sprintf("%v", p.Opts),
		})
	}

	return disks, nil
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

// Sends disk data to the server
func PostDiskData(a *agent.Agent, data []AgentDisk) error {
	// Wrap data inside "name" and "data" keys
	payload, err := json.Marshal(map[string]interface{}{
		"name": "get-disk",
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
		log.Printf("Failed to upload disk data, status code: %d\n", resp.StatusCode)
		return fmt.Errorf("failed to upload disk data, status code: %d", resp.StatusCode)
	}

	return nil
}
