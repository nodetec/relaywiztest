package cmd

import (
	"fmt"
	"github.com/nodetec/relaywiz/pkg/manager"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the installation and configuration of the nostr relay...")
		apt.InstallPackages()
		// Add other steps here for nginx config, firewall setup, etc.
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

