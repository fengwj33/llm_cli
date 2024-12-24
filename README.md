# llmcli

A command line interface for interacting with various LLM APIs (OpenAI, ChatGLM) with support for assistants and chat history.

## Features

- Support for multiple LLM providers (ChatGLM, OpenAI)
- Assistant system with customizable prompts
- Chat history management
- Pipe support for processing file content
- Markdown rendering for responses
- Configurable chat context window

## Installation

### From Source

Requirements:
- Go 1.19 or higher
- GCC (for SQLite support)

Build Steps:
1. git clone https://github.com/yourusername/llmcli.git
2. cd llmcli
3. go mod download
4. go build -o llmcli
5. Optional: sudo mv llmcli /usr/local/bin/

## Configuration

Before first run, use `llmcli -c` to create and edit the configuration file at ~/.config/llm_cli/config.json. This will allow you to add your API keys and configure models and assistants.
Example configuration structure:

1. Default Model:
   - Specify which model to use when no model is explicitly selected
   - Example: "default": "gpt4"

2. Models:
   - Configure direct API access to language models
   - Each model entry requires:
     - API: provider name (e.g., "openai", "chatglm")
     - Model: specific model identifier (e.g., "gpt-4", "chatglm-6b")
     - API_KEY: authentication key for the API
   - Used with -m flag for one-off queries without context

3. Assistants:
   - Define specialized chat interfaces with persistent context
   - Each assistant entry includes:
     - model: which model to use (must match a configured model name)
     - prompt: system prompt that defines assistant's behavior
     - chatContextWindow: number of previous exchanges to include
   - Used with -a flag for contextual conversations
   - Maintains chat history for continuous dialogue

Example use case:
- Models: Direct questions like "llmcli -m gpt4 'what is 2+2?'"
- Assistants: Complex tasks like "llmcli -a code_review 'review this function'"
  where the assistant remembers context and follows a specific prompt

## Usage

### Basic Commands
- llmcli "hello world" - Use default assistant
- llmcli -m gpt4 "what is golang" - Use specific model
- llmcli -a coding "tell me about channels" - Use specific assistant
- llmcli -c - Edit configuration
- echo "some text" | llmcli - Process text from pipe

### Chat History Commands
- llmcli -h assistant_name - Show chat history
- llmcli -h assistant_name 5 - Show last 5 messages
- llmcli --clear assistant_name - Clear chat history

## Examples

### Using Default Assistant
```shell
llmcli "explain what is golang"
```

### Using Specific Model
```shell
llmcli -m chatglm "tell me a story"
```

### Using Assistant with Pipe
```shell
cat code.go | llmcli -a code_review
```

### Using History
```shell
llmcli -h code_review 10
```

### Using Aliases for Quick Access
Add aliases to your shell configuration file (.bashrc, .zshrc, etc.) for faster access:

For example, add these lines to your .zshrc:
```shell
alias a1="llmcli -a assistant1"    # Quick access to assistant1
alias a2="llmcli -a code_review"   # Quick access to code review assistant
alias a3="llmcli -a translator"    # Quick access to translator assistant
```

Then you can simply use:
```shell
a1 "hello"              # Chat with assistant1
a2 "review this code"   # Start code review
a3 "translate to zh"    # Use translator
```

## Project Structure

- config/ - Configuration management
- llm/ - Core LLM functionality
  - api/ - API providers implementation
  - assistant.go - Assistant functionality
  - llm.go - Main LLM interface
- utils/ - Utility functions
- tests/ - Unit tests
- main.go - CLI entry point

## Development

To add a new LLM provider:
1. Create new provider file in llm/api/
2. Implement the LLMProvider interface
3. Add provider to the Providers map in api.go

## License

Licensed under the Apache License, Version 2.0. See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
