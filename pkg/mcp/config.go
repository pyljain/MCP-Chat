package mcp

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Name    string            `json:"name"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

func LoadConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		return &Config{}
	}

	configPath := fmt.Sprintf("%s/.merlin/mcp.json", homeDir)
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return &Config{}
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		return &Config{}
	}

	return &config
}
