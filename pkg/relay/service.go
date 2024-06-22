package relay

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const serviceTemplate = `[Unit]
Description=Nostr Relay Pyramid
After=network.target

[Service]
Type=simple
User=nostr
WorkingDirectory=/home/nostr
Environment="DOMAIN={{.Domain}}"
Environment="RELAY_NAME=nostr-relay-pyramid"
Environment="RELAY_PUBKEY=asdf"
ExecStart=/usr/local/bin/nostr-relay-pyramid
Restart=on-failure

[Install]
WantedBy=multi-user.target
`

const serviceFilePath = "/etc/systemd/system/nostr-relay-pyramid.service"

func userExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	err := cmd.Run()
	return err == nil
}

func SetupRelayService(domain string) {
	if _, err := os.Stat(serviceFilePath); err == nil {
		fmt.Printf("Service file already exists at %s.\n", serviceFilePath)
		return
	}

	if !userExists("nostr") {
		fmt.Println("Creating user for nostr relay service...")
		err := exec.Command("adduser", "--disabled-login", "--gecos", "", "nostr").Run()
		if err != nil {
			log.Fatalf("Error creating user: %v", err)
		}
	} else {
		fmt.Println("User 'nostr' already exists.")
	}

	fmt.Println("Creating data directory...")
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	fmt.Println("Setting ownership of the data directory...")
	err = os.Chown(dataDir, os.Getuid(), os.Getgid())
	if err != nil {
		log.Fatalf("Error setting ownership of the data directory: %v", err)
	}

	fmt.Println("Creating systemd service file...")
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		log.Fatalf("Error creating service file: %v", err)
	}
	defer serviceFile.Close()

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		log.Fatalf("Error parsing service template: %v", err)
	}

	err = tmpl.Execute(serviceFile, struct{ Domain string }{Domain: domain})
	if err != nil {
		log.Fatalf("Error executing service template: %v", err)
	}

	fmt.Println("Reloading systemd daemon...")
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		log.Fatalf("Error reloading systemd daemon: %v", err)
	}

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

