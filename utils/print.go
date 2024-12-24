package utils

import "fmt"

const (
	RedColor   = "\033[31m"
	ResetColor = "\033[0m"
)

// PrintError prints error messages in red color
func PrintError(format string, a ...interface{}) {
	fmt.Printf(RedColor+format+ResetColor+"\n", a...)
} 