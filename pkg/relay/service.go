package relay

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Template for the systemd service file
const serviceTemplate = `[Unit]
Description=Nostr Relay Pyramid
After=network.target

[Service]
Type=simple
User=nostr
WorkingDirectory=/home/nostr
Environment="DOMAIN=test"
Environment="RELAY_NAME=nostr-relay-pyramid"
Environment="RELAY_PUBKEY=asdf"
ExecStart=/usr/local/bin/nostr-relay-pyramid
Restart=on-failure

[Install]
WantedBy=multi-user.target
`

// Path for the systemd service file
const serviceFilePath = "/etc/systemd/system/nostr-relay-pyramid.service"

// Function to check if a user exists
func userExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	err := cmd.Run()
	return err == nil
}

// Function to create the systemd service file and start the service
func SetupRelayService() {
	// Check if the service file already exists
	if _, err := os.Stat(serviceFilePath); err == nil {
		fmt.Printf("Service file already exists at %s.\n", serviceFilePath)
		return
	}

	// Check if the user already exists
	if !userExists("nostr") {
		// Create a user for the nostr relay service
		fmt.Println("Creating user for nostr relay service...")
		err := exec.Command("adduser", "--disabled-login", "--gecos", "", "nostr").Run()
		if err != nil {
			log.Fatalf("Error creating user: %v", err)
		}
	} else {
		fmt.Println("User 'nostr' already exists.")
	}

	// Set ownership of the data directory
	fmt.Println("Setting ownership of the data directory...")
	err := os.Chown(dataDir, os.Getuid(), os.Getgid())
	if err != nil {
		log.Fatalf("Error setting ownership of the data directory: %v", err)
	}

	// Create the systemd service file
	fmt.Println("Creating systemd service file...")
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		log.Fatalf("Error creating service file: %v", err)
	}
	defer serviceFile.Close()

	_, err = serviceFile.WriteString(serviceTemplate)
	if err != nil {
		log.Fatalf("Error writing to service file: %v", err)
	}

	// Reload systemd to apply the new service
	fmt.Println("Reloading systemd daemon...")
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		log.Fatalf("Error reloading systemd daemon: %v", err)
	}

	// Enable and start the nostr relay service
	fmt.Println("Enabling and starting nostr relay service...")
	err = exec.Command("systemctl", "enable", "nostr-relay-pyramid").Run()
	if err != nil {
		log.Fatalf("Error enabling nostr relay service: %v", err)
	}

	err = exec.Command("systemctl", "start", "nostr-relay-pyramid").Run()
	if err != nil {
		log.Fatalf("Error starting nostr relay service: %v", err)
	}

	fmt.Println("Nostr relay service setup completed.")
}

