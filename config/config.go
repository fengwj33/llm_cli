package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"llm_cli/utils"
)

const (
	ConfigDir  = ".config/llm_cli"
	ConfigFile = "config.json"
)

// HandleConfig opens the configuration file in vi editor
func HandleConfig() {
	configPath := GetConfigPath()
	
	if err := BackupConfig(configPath); err != nil {
		utils.PrintError("Error backing up config: %v", err)
		return
	}

	if err := openInVi(configPath); err != nil {
		utils.PrintError("Error opening vi: %v", err)
		return
	}

	if err := ValidateConfig(configPath); err != nil {
		utils.PrintError("Invalid JSON configuration: %v", err)
		utils.PrintError("Reverting to previous version...")
		if restoreErr := RestoreConfig(configPath); restoreErr != nil {
			utils.PrintError("Error restoring config: %v", restoreErr)
		}
	}
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configDirPath := filepath.Join(home, ConfigDir)
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return ""
	}
	
	configPath := filepath.Join(configDirPath, ConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.WriteFile(configPath, []byte(DefaultConfigTemplate), 0644)
	}
	return configPath
}

func BackupConfig(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath+".backup", content, 0644)
}

func RestoreConfig(configPath string) error {
	content, err := os.ReadFile(configPath + ".backup")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, content, 0644)
}

func ValidateConfig(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	
	var jsonData interface{}
	return json.Unmarshal(content, &jsonData)
}

func openInVi(configPath string) error {
	cmd := exec.Command("vi", configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
} 