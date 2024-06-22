package cmd

import (
	"fmt"

	"github.com/pterm/pterm"

	"github.com/nodetec/relaywiz/pkg/manager"
	"github.com/nodetec/relaywiz/pkg/network"
	"github.com/nodetec/relaywiz/pkg/relay"
	"github.com/nodetec/relaywiz/pkg/ui"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {

		ui.Greet()

		relayDomain, _ := pterm.DefaultInteractiveTextInput.Show("Relay domain name")
		email, _ := pterm.DefaultInteractiveTextInput.Show("Email address")

		fmt.Printf("Starting the installation and configuration of the nostr relay with domain: %s and email: %s...\n", relayDomain, email)

		// Step 1: Install necessary packages using APT
		manager.AptInstallPackages()

		// Step 2: Configure the firewall
		network.ConfigureFirewall()

		// Step 3: Configure Nginx for HTTP
		network.ConfigureNginxHttp(relayDomain)

		// Step 4: Get SSL certificates
		network.GetCertificates(relayDomain, email)

		// Step 5: Configure Nginx for HTTPS
		network.ConfigureNginxHttps(relayDomain)

		// Step 6: Download and install the relay binary
		relay.InstallRelayBinary()

		// Step 7: Set up the relay service
		relay.SetupRelayService(relayDomain)

		// Add other steps here for starting the relay service, etc.
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
