package network

import (
	"fmt"
	"log"
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
		fmt.Printf("SSL certificates for %s already exist.\n", domainName)
		return
	}

	fmt.Println("Creating necessary directories for Certbot...")
	err := os.MkdirAll(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", dirName), 0755)
	if err != nil {
		log.Fatalf("Error creating directories for Certbot: %v", err)
	}

	fmt.Printf("Obtaining SSL certificates for %s using Certbot...\n", domainName)
	cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("/var/www/%s", dirName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Certbot failed to obtain the certificate for %s: %v", domainName, err)
	}

	fmt.Printf("SSL certificates for %s obtained successfully.\n", domainName)
}

