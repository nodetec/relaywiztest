package cmd

import (
	"fmt"
	"github.com/nodetec/relaywiz/pkg/manager"
	"github.com/nodetec/relaywiz/pkg/network"
	"github.com/nodetec/relaywiz/pkg/relay"

	"github.com/spf13/cobra"
)

// Define default domain name and email
const domainName = "relay.testrelay.xyz"
const email = "chris.machine@pm.me"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the installation and configuration of the nostr relay...")

		// Step 1: Install necessary packages using APT
		manager.AptInstallPackages()

		// Step 2: Configure the firewall
		network.ConfigureFirewall()

		// Step 3: Configure Nginx for HTTP
		network.ConfigureNginxHttp(domainName)

		// Step 4: Get SSL certificates
		network.GetCertificates(domainName, email)

		// Step 5: Configure Nginx for HTTPS
		network.ConfigureNginxHttps(domainName)

		// Step 6: Download and install the relay binary
		relay.InstallRelayBinary()

		// Step 7: Set up the relay service
		relay.SetupRelayService()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
