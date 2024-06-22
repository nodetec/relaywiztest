package utils

import (
	"fmt"

	"os"
	"strings"
)

// Function to print info messages
func PrintInfo(message string) {
	fmt.Println(message)
}

// Function to print success messages
func PrintSuccess(message string) {
	fmt.Println(message)
}

// Function to print error messages
func PrintError(message string) {
	fmt.Println(message)
}

// Function to check if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && !info.IsDir()
}

// Function to extract the directory name from the domain
func GetDirectoryName(domainName string) string {
	domainParts := strings.Split(domainName, ".")
	if len(domainParts) > 2 {
		return domainParts[1]
	}
	return domainParts[0]
}
