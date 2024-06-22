package cmd

import (
	"fmt"
	"github.com/nodetec/relaywiz/pkg/manager"
	"github.com/nodetec/relaywiz/pkg/network"

	"github.com/spf13/cobra"
)

// Define a domain name (this can be replaced with a configuration or environment variable)
const domainName = "yourdomain.com"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the installation and configuration of the nostr relay...")

		// Step 1: Install necessary packages using APT
		manager.AptInstallPackages()

		// Step 2: Configure Nginx for HTTP
		network.ConfigureNginxHttp(domainName)

		// Add other steps here for firewall setup, SSL certificates, etc.
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

