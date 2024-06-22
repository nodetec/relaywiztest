package network

import (
	"log"
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to configure the firewall
func ConfigureFirewall() {
	spinner, _ := pterm.DefaultSpinner.Start("Configuring firewall to allow HTTP (port 80) and HTTPS (port 443) traffic...")
	exec.Command("ufw", "enable").Run()
	// spinner.UpdateText("Firewall enabled successfully.")

	// Allow HTTP and HTTPS traffic
	err := exec.Command("ufw", "allow", "Nginx Full").Run()
	// spinner.UpdateText("Allowing Nginx Full...")
	if err != nil {
		log.Fatalf("Error allowing Nginx Full: %v", err)
	}

	// Reload the firewall to apply the changes
	err = exec.Command("ufw", "reload").Run()
	// spinner.UpdateText("Reloading firewall...")
	if err != nil {
		log.Fatalf("Error reloading firewall: %v", err)
	}

	spinner.Success("Firewall configured successfully.")
}
