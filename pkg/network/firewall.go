package network

import (
	"fmt"
	"os/exec"
	"github.com/nodetec/relaywiz/pkg/utils"
)

// Function to configure the firewall
func ConfigureFirewall() {
	utils.PrintInfo("Configuring firewall to allow HTTP (port 80) and HTTPS (port 443) traffic...")

	// Allow HTTP and HTTPS traffic
	err := exec.Command("ufw", "allow", "Nginx Full").Run()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Error allowing Nginx Full: %v", err))
		return
	}

	// Reload the firewall to apply the changes
	err = exec.Command("ufw", "reload").Run()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Error reloading firewall: %v", err))
		return
	}

	// Show the current firewall status
	err = exec.Command("ufw", "status", "verbose").Run()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Error showing firewall status: %v", err))
		return
	}

	utils.PrintSuccess("Firewall configuration completed successfully.")
}

