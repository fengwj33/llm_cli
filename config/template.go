package config

const DefaultConfigTemplate = `{
	"default": "chatglm",
	"models": {
		"gpt-4o": {
			"API": "OpenAI",
			"Model": "gpt-4o",
			"API_KEY": "your-api-key-here"
		},
		"chatglm": {
			"API": "ChatGLM",
			"Model": "glm-4v-flash",
			"API_KEY": "your-api-key-here"
		}
	},
	"assistants": {
		"assistant": {
			"model": "chatglm",
			"prompt": "You are a helpful assistant.",
			"chatContextWindow":5
		},
		"code_reviewer": {
			"model": "chatglm",
			"prompt": "You are a helpful code reviewer and you will review the code and provide feedback.",
			"chatContextWindow":5
		}
	}
}
` 