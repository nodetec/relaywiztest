package network

import (
	"fmt"
	"log"
	"os/exec"
)

// Function to configure the firewall
func ConfigureFirewall() {
	fmt.Println("Configuring firewall to allow HTTP (port 80) and HTTPS (port 443) traffic...")

	exec.Command("ufw", "enable").Run()

	// Allow HTTP and HTTPS traffic
	err := exec.Command("ufw", "allow", "Nginx Full").Run()
	if err != nil {
		log.Fatalf("Error allowing Nginx Full: %v", err)
	}

	// Reload the firewall to apply the changes
	err = exec.Command("ufw", "reload").Run()
	if err != nil {
		log.Fatalf("Error reloading firewall: %v", err)
	}

	// Show the current firewall status
	err = exec.Command("ufw", "status", "verbose").Run()
	if err != nil {
		log.Fatalf("Error showing firewall status: %v", err)
	}

	fmt.Println("Firewall configuration completed successfully.")
}
