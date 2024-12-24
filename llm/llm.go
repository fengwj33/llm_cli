package llm

import (
	"fmt"
	"llm_cli/config"
	"llm_cli/llm/api"
)

// Request represents a request to an LLM model
type Request struct {
	Model    string
	Messages []api.Message
}

// Call sends a request to the specified LLM model and returns its response
func Call(modelName string, messages []api.Message) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}

	model, exists := cfg.Models[modelName]
	if !exists {
		if cfg.Default != "" {
			model, exists = cfg.Models[cfg.Default]
			if !exists {
				return "", fmt.Errorf("default model '%s' not found in config", cfg.Default)
			}
		} else {
			return "", fmt.Errorf("model '%s' not found in config", modelName)
		}
	}

	provider, exists := api.Providers[model.API]
	if !exists {
		return "", fmt.Errorf("unsupported API provider: %s", model.API)
	}

	return provider.Call(model.Model, messages, model.API_KEY)
}

// SimpleCall is a helper function for simple single-message calls
func SimpleCall(modelName string, input string) (string, error) {
	messages := []api.Message{
		{Role: "user", Content: input},
	}
	return Call(modelName, messages)
} 