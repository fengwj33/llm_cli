package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"llm_cli/config"
	"llm_cli/llm"
	"llm_cli/utils"

	"github.com/charmbracelet/glamour"
)

const (
	usageTemplate = `Usage:
  llmcli <text>                        - Use default assistant with text
  llmcli -c, --config                 - Edit configuration file
  llmcli -m, --model <name> <text>    - Call specific model with text
  llmcli -a, --assistant <name> <text> - Call specific assistant with text
  llmcli -h, --history <name> [n]     - Show chat history for assistant (last n messages)
  llmcli --clear <name>               - Clear chat history for assistant`
)

func getInput() string {
	// Check if there's input from pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read from pipe
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
			return ""
		}
		return string(bytes)
	}
	
	return "" // Remove command line argument handling from here
}

func main() {
	if len(os.Args) == 1 {
		// Check for pipe input
		if input := getInput(); input != "" {
			handleAssistantCall("", input)
			return
		}
		showUsage()
		return
	}

	switch os.Args[1] {
	case "-c", "--config":
		config.HandleConfig()
	case "-d", "--debug":
		debug()
	case "-h", "--history":
		handleHistory(os.Args[2:])
	case "--clear":
		handleClearHistory(os.Args[2:])
	case "-m", "--model":
		if len(os.Args) < 3 {
			fmt.Println("Error: Model name required")
			fmt.Println("Usage: llmcli -m <model_name> [input_text]")
			return
		}
		modelName := os.Args[2]
		input := getInput() // Check pipe input first
		if input == "" && len(os.Args) > 3 {
			input = strings.Join(os.Args[3:], " ") // Use remaining args as input
		}
		if input == "" {
			fmt.Println("Error: No input provided")
			return
		}
		handleModelCall(modelName, input)
	case "-a", "--assistant":
		if len(os.Args) < 3 {
			fmt.Println("Error: Assistant name required")
			fmt.Println("Usage: llmcli -a <assistant_name> [input_text]")
			return
		}
		assistantName := os.Args[2]
		input := getInput() // Check pipe input first
		if input == "" && len(os.Args) > 3 {
			input = strings.Join(os.Args[3:], " ") // Use remaining args as input
		}
		if input == "" {
			fmt.Println("Error: No input provided")
			return
		}
		handleAssistantCall(assistantName, input)
	default:
		input := getInput() // Check pipe input first
		if input == "" {
			input = strings.Join(os.Args[1:], " ") // Use all args as input
		}
		handleAssistantCall("", input)
	}
}

func showUsage() {
	fmt.Println(usageTemplate)
}

func renderResponse(response string) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Printf("Error initializing renderer: %v\n", err)
		fmt.Println(response) // Fallback to plain text
		return
	}

	out, err := r.Render(response)
	if err != nil {
		fmt.Printf("Error rendering markdown: %v\n", err)
		fmt.Println(response) // Fallback to plain text
		return
	}

	fmt.Print(out)
}

func handleModelCall(modelName, input string) {
	response, err := llm.SimpleCall(modelName, input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	renderResponse(response)
}

func handleAssistantCall(assistantName string, input string) {
	var response string
	var err error

	if assistantName == "" {
		response, err = llm.SimpleAssistantCall(input)
	} else {
		response, err = llm.AssistantCall(assistantName, input)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	renderResponse(response)
}

func handleHistory(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Assistant name required")
		fmt.Println("Usage: echo -h <assistant_name> [number_of_messages]")
		return
	}

	assistantName := args[0]
	limit := 10 // default limit
	if len(args) > 1 {
		if n, err := strconv.Atoi(args[1]); err == nil && n > 0 {
			limit = n
		}
	}

	history, err := utils.NewHistory()
	if err != nil {
		fmt.Printf("Error initializing history: %v\n", err)
		return
	}
	defer history.Close()

	records, err := history.Fetch(assistantName, limit)
	if err != nil {
		fmt.Printf("Error fetching history: %v\n", err)
		return
	}

	if len(records) == 0 {
		fmt.Printf("No chat history found for assistant '%s'\n", assistantName)
		return
	}

	fmt.Printf("\nChat history for assistant '%s' (last %d messages):\n", assistantName, limit)
	fmt.Println("----------------------------------------")

	// Print in chronological order (oldest first)
	for _, record := range records {
		roleColor := "\033[36m" // cyan for user
		if record.Role == "assistant" {
			roleColor = "\033[32m" // green for assistant
		}
		fmt.Printf("%s%s\033[0m: %s\n\n", roleColor, record.Role, record.Content)
	}
}

func handleClearHistory(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Assistant name required")
		fmt.Println("Usage: echo --clear <assistant_name>")
		return
	}

	assistantName := args[0]
	history, err := utils.NewHistory()
	if err != nil {
		fmt.Printf("Error initializing history: %v\n", err)
		return
	}
	defer history.Close()

	if err := history.Clear(assistantName); err != nil {
		fmt.Printf("Error clearing history: %v\n", err)
		return
	}

	fmt.Printf("Successfully cleared chat history for assistant '%s'\n", assistantName)
}

func debug() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		return
	}

	// Test models
	fmt.Println("\n=== Testing Models ===")
	for name, model := range cfg.Models {
		fmt.Printf("\nModel: %s\n", name)
		fmt.Printf("  API: %s\n", model.API)
		fmt.Printf("  Model: %s\n", model.Model)
		fmt.Printf("  API_KEY: %s\n", model.API_KEY)

		response, err := llm.SimpleCall(name, "This is a test message")
		if err != nil {
			fmt.Printf("  Test call error: %v\n", err)
		} else {
			fmt.Printf("  Test response: %s\n", response)
		}
	}

	// Test assistants
	fmt.Println("\n=== Testing Assistants ===")
	for name, assistant := range cfg.Assistants {
		fmt.Printf("\nAssistant: %s\n", name)
		fmt.Printf("  Model: %s\n", assistant.Model)
		fmt.Printf("  Prompt: %s\n", assistant.Prompt)
		fmt.Printf("  ChatContextWindow: %d\n", assistant.ChatContextWindow)

		response, err := llm.AssistantCall(name, "This is a test message")
		if err != nil {
			fmt.Printf("  Test call error: %v\n", err)
		} else {
			fmt.Printf("  Test response: %s\n", response)
		}
	}
}