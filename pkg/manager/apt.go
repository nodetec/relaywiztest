package manager

import (
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to check if a command exists
func commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Function to install necessary packages
func AptInstallPackages() {
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Updating and installing packages...")
	exec.Command("apt", "update", "-qq").Run()

	// Check if nginx is installed, install if not
	if commandExists("nginx") {
	} else {
		exec.Command("apt", "install", "-y", "-qq", "nginx").Run()
	}

	// Check if Certbot is installed, install if not
	if commandExists("certbot") {
	} else {
		exec.Command("apt", "install", "-y", "-qq", "certbot", "python3-certbot-nginx").Run()
	}

	// Check if ufw is installed, install if not
	if commandExists("ufw") {
	} else {
		exec.Command("apt", "install", "-y", "-qq", "ufw").Run()
		exec.Command("ufw", "enable").Run()
	}

	spinnerInfo.Success()
}

// Function to check if a package is installed
func isPackageInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	err := cmd.Run()
	return err == nil
}
