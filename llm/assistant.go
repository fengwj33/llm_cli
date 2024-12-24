package llm

import (
	"fmt"
	"llm_cli/config"
	"llm_cli/llm/api"
	"llm_cli/utils"
)

// AssistantCall sends a request using a configured assistant
func AssistantCall(assistantName string, input string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}

	// Get assistant config
	assistant, exists := cfg.Assistants[assistantName]
	if !exists {
		return "", fmt.Errorf("assistant '%s' not found in config", assistantName)
	}

	// Initialize history
	history, err := utils.NewHistory()
	if err != nil {
		return "", fmt.Errorf("failed to initialize history: %v", err)
	}
	defer history.Close()

	// Get recent chat context
	records, err := history.Fetch(assistantName, assistant.ChatContextWindow*2) // *2 because each turn has 2 messages
	if err != nil {
		return "", fmt.Errorf("failed to fetch history: %v", err)
	}

	// Build messages array starting with system prompt
	messages := []api.Message{
		{
			Role:    "system",
			Content: assistant.Prompt,
		},
	}

	// Add context messages in chronological order
	// No need to reverse since Fetch now returns them in chronological order
	for _, record := range records {
		messages = append(messages, api.Message{
			Role:    record.Role,
			Content: record.Content,
		})
	}

	// Add current user input
	messages = append(messages, api.Message{
		Role:    "user",
		Content: input,
	})

	// Call the model
	response, err := Call(assistant.Model, messages)
	if err != nil {
		return "", fmt.Errorf("model call failed: %v", err)
	}

	// Store the conversation in history
	if err := history.Push(assistantName, "user", input); err != nil {
		return "", fmt.Errorf("failed to store user message: %v", err)
	}
	if err := history.Push(assistantName, "assistant", response); err != nil {
		return "", fmt.Errorf("failed to store assistant response: %v", err)
	}

	return response, nil
}

// SimpleAssistantCall uses the default assistant if none specified
func SimpleAssistantCall(input string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}

	// Use first assistant as default if none specified
	var defaultAssistant string
	for name := range cfg.Assistants {
		defaultAssistant = name
		break
	}

	if defaultAssistant == "" {
		return "", fmt.Errorf("no assistants configured")
	}

	return AssistantCall(defaultAssistant, input)
} 