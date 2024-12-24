package api

// LLMProvider defines the interface for LLM providers
type LLMProvider interface {
	Call(model string, messages []Message, apiKey string) (string, error)
}

// BaseProvider implements common functionality
type BaseProvider struct {
	Name string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Provider map to store available providers
var Providers = map[string]LLMProvider{
	"OpenAI":   &OpenAIProvider{BaseProvider{Name: "OpenAI"}},
	"ChatGLM":  &ChatGLMProvider{BaseProvider{Name: "ChatGLM"}},
} 