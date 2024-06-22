package network

import (
	"fmt"
  "os"
	"os/exec"
	"github.com/nodetec/relaywiz/pkg/utils"
)

// Function to get SSL certificates using Certbot
func GetCertificates(domainName, email string) {
	dirName := utils.GetDirectoryName(domainName)

	// Check if certificates already exist
	if utils.FileExists(fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domainName)) &&
		utils.FileExists(fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domainName)) {
		utils.PrintSuccess(fmt.Sprintf("SSL certificates for %s already exist.", domainName))
		return
	}

	utils.PrintInfo("Creating necessary directories for Certbot...")
	err := os.MkdirAll(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", dirName), 0755)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Error creating directories for Certbot: %v", err))
		return
	}

	utils.PrintInfo(fmt.Sprintf("Obtaining SSL certificates for %s using Certbot...", domainName))
	cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("/var/www/%s", dirName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
	err = cmd.Run()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Certbot failed to obtain the certificate for %s: %v", domainName, err))
		return
	}

	utils.PrintSuccess(fmt.Sprintf("SSL certificates for %s obtained successfully.", domainName))
}

