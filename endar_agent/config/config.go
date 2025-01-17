package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// Parse command-line flags and return them as a map
func ParseFlags() map[string]interface{} {
	flags := make(map[string]interface{})

	// Configurable parameters
	flags["server"] = flag.String("server", "", "Server URL")
	flags["key"] = flag.String("key", "", "Registration Token")
	flags["config"] = flag.String("config", "config.ini", "Path to configuration file")

	// Service command flags
	flags["install"] = flag.Bool("install", false, "Install the service")
	flags["uninstall"] = flag.Bool("uninstall", false, "Uninstall the service")
	flags["start"] = flag.Bool("start", false, "Start the service")
	flags["stop"] = flag.Bool("stop", false, "Stop the service")
	flags["restart"] = flag.Bool("restart", false, "Restart the service")
	flags["status"] = flag.Bool("status", false, "Check the service status")

	flag.Parse()

	// Convert pointer values to actual values in the map
	for key, value := range flags {
		switch v := value.(type) {
		case *string:
			flags[key] = *v
		case *bool:
			flags[key] = *v
		}
	}

	return flags
}

// Load configuration from file
func LoadConfig(configPath string) (*ini.File, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config.ini: %v", err)
	}
	return cfg, nil
}

// Save configuration to file
func SaveConfig(configPath string, cfg *ini.File) error {
	return cfg.SaveTo(configPath)
}

// Save default configuration
func SaveDefaultConfig(configPath, serverURL, registrationToken string) error {
	cfg := ini.Empty()
	cfg.Section("configuration").Key("server").SetValue(serverURL)
	cfg.Section("configuration").Key("key").SetValue(registrationToken)
	cfg.Section("configuration").Key("interval").SetValue("300")
	cfg.Section("configuration").Key("level").SetValue("DEBUG")
	cfg.Section("configuration").Key("print_config").SetValue("False")
	cfg.Section("configuration").Key("aid").SetValue("")
	cfg.Section("configuration").Key("registered").SetValue("False")
	return SaveConfig(configPath, cfg)
}

// Initialize configuration and return settings
func InitializeConfig(flags map[string]interface{}) (string, string, int, string, string, bool, string, bool) {
	configPath := flags["config"].(string)
	serverURL := flags["server"].(string)
	registrationToken := flags["key"].(string)

	// Load or create default config
	var cfg *ini.File
	var err error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("No config file found, creating new one...")
		if err := SaveDefaultConfig(configPath, serverURL, registrationToken); err != nil {
			log.Fatalf("Failed to save default config: %v", err)
		}
	}

	cfg, err = LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Read config values
	server := cfg.Section("configuration").Key("server").String()
	key := cfg.Section("configuration").Key("key").String()
	interval, _ := cfg.Section("configuration").Key("interval").Int()
	debugLevel := cfg.Section("configuration").Key("level").String()
	configFilePath := cfg.Section("configuration").Key("config_path").String()
	printConfig, _ := cfg.Section("configuration").Key("print_config").Bool()
	aid := cfg.Section("configuration").Key("aid").String()
	registered, _ := cfg.Section("configuration").Key("registered").Bool()

	// Override config with flags if provided
	if serverURL != "" {
		server = serverURL
		cfg.Section("configuration").Key("server").SetValue(serverURL)
	}
	if registrationToken != "" {
		key = registrationToken
		cfg.Section("configuration").Key("key").SetValue(registrationToken)
	}
	if configPath != "" {
		configFilePath = configPath
		cfg.Section("configuration").Key("config_path").SetValue(configPath)
	}

	// Save updated config
	SaveConfig(configPath, cfg)

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
