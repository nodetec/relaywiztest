package apt

import (
	"fmt"
	"os/exec"
)

// Function to print info messages
func printInfo(message string) {
	fmt.Println(message)
}

// Function to print success messages
func printSuccess(message string) {
	fmt.Println(message)
}

// Function to check if a command exists
func commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Function to install necessary packages
func InstallPackages() {
	printInfo("Updating package lists silently...")
	exec.Command("apt", "update", "-qq").Run()

	// Check if nginx is installed, install if not
	if commandExists("nginx") {
		printSuccess("nginx is already installed.")
	} else {
		printInfo("Installing nginx...")
		exec.Command("apt", "install", "-y", "-qq", "nginx").Run()
	}

	// Check if Certbot is installed, install if not
	if commandExists("certbot") {
		printSuccess("Certbot is already installed.")
	} else {
		printInfo("Installing Certbot and dependencies...")
		exec.Command("apt", "install", "-y", "-qq", "certbot", "python3-certbot-nginx").Run()
	}

	// Check if ufw is installed, install if not
	if commandExists("ufw") {
		printSuccess("ufw is already installed.")
	} else {
		printInfo("Installing ufw...")
		exec.Command("apt", "install", "-y", "-qq", "ufw").Run()
		printInfo("Enabling ufw...")
		exec.Command("ufw", "enable").Run()
	}

	// Check if the nostr relay dependencies are installed, install if not
	checkAndInstallDependencies()
}

// Function to check and install nostr relay dependencies
func checkAndInstallDependencies() {
	dependencies := []string{
		"build-essential",
		"cmake",
		"protobuf-compiler",
		"pkg-config",
		"libssl-dev",
		"git",
	}

	allInstalled := true
	for _, dep := range dependencies {
		if !isPackageInstalled(dep) {
			allInstalled = false
			break
		}
	}

	if allInstalled {
		printSuccess("All nostr relay dependencies are already installed.")
	} else {
		printInfo("Installing dependencies for nostr relay...")
		args := append([]string{"install", "-y", "-qq"}, dependencies...)
		exec.Command("apt", args...).Run()
	}
}

// Function to check if a package is installed
func isPackageInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	err := cmd.Run()
	return err == nil
}
