package config

import (
	"encoding/json"
	"os"
	"sync"
)

// ModelConfig represents the configuration for a single model
type ModelConfig struct {
	API     string `json:"API"`
	Model   string `json:"Model"`
	API_KEY string `json:"API_KEY"`
}

// AssistantConfig represents the configuration for an assistant
type AssistantConfig struct {
	Model             string `json:"model"`
	Prompt            string `json:"prompt"`
	ChatContextWindow int    `json:"chatContextWindow"`
}

// Config represents the root configuration structure
type Config struct {
	Default    string                     `json:"default"`
	Models     map[string]ModelConfig     `json:"models"`
	Assistants map[string]AssistantConfig `json:"assistants"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

// GetConfig returns the singleton instance of Config
func GetConfig() (*Config, error) {
	once.Do(func() {
		configPath := GetConfigPath()
		var err error
		instance, err = LoadConfig(configPath)
		if err != nil {
			instance = &Config{
				Models:     make(map[string]ModelConfig),
				Assistants: make(map[string]AssistantConfig),
			}
		}
	})
	return instance, nil
}

// LoadConfig loads and parses the configuration file
func LoadConfig(configPath string) (*Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config, configPath string) error {
	mu.Lock()
	defer mu.Unlock()
	
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
} 