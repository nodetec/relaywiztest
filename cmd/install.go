package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nodetec/relaywiz/pkg/manager"
	"github.com/nodetec/relaywiz/pkg/network"
	"github.com/nodetec/relaywiz/pkg/relay"
	"github.com/spf13/cobra"
)

type model struct {
	input string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		case "ctrl+c", "esc":
			return m, tea.Quit
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Enter the relay domain name: %s", m.input)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new Bubble Tea program
		p := tea.NewProgram(model{})

		// Run the Bubble Tea program to get user input
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Extract the user input
		relayDomain := m.(model).input
    const email = "chris.machine@pm.me"

		fmt.Printf("Starting the installation and configuration of the nostr relay with domain: %s...\n", relayDomain)

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
		relay.SetupRelayService()

		// Add other steps here for starting the relay service, etc.
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

