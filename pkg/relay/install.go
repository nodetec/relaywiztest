package relay

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// URL of the binary to download
const downloadURL = "https://github.com/github-tijlxyz/khatru-pyramid/releases/download/v0.0.5/khatru-pyramid-v0.0.5-linux-amd64"

// Name of the binary after downloading
const binaryName = "nostr-relay-pyramid"

// Destination directory for the binary
const destDir = "/usr/local/bin"

// Function to download and make the binary executable
func InstallRelayBinary() {
	// Determine the file name from the URL
	tempFileName := filepath.Base(downloadURL)

	// Create the temporary file
	out, err := os.Create(tempFileName)
	if err != nil {
		log.Fatalf("Error creating temporary file: %v", err)
	}
	defer out.Close()

	// Download the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad status: %s", resp.Status)
	}

	// Write the body to the temporary file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Error writing to temporary file: %v", err)
	}

	// Define the final destination path
	destPath := filepath.Join(destDir, binaryName)

	// Move the file to the destination directory
	err = os.Rename(tempFileName, destPath)
	if err != nil {
		log.Fatalf("Error moving file to /usr/local/bin: %v", err)
	}

	// Make the file executable
	err = os.Chmod(destPath, 0755)
	if err != nil {
		log.Fatalf("Error making file executable: %v", err)
	}

	fmt.Printf("%s downloaded, moved to %s, and made executable successfully.\n", tempFileName, destPath)
}
