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
	inputs    []string
	cursor    int
	input     string
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.inputs = append(m.inputs, m.input)
			m.input = ""
			if m.cursor < 1 {
				m.cursor++
			} else {
				return m, tea.Quit
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Exiting the program. Goodbye!\n"
	}
	switch m.cursor {
	case 0:
		return fmt.Sprintf("Enter the relay domain name: %s", m.input)
	case 1:
		return fmt.Sprintf("Enter the email address: %s", m.input)
	default:
		return "Thank you! Proceeding with installation..."
	}
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new Bubble Tea program
		p := tea.NewProgram(model{inputs: make([]string, 0)})

		// Run the Bubble Tea program to get user input
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Check if the program was exited with ctrl+c or esc
		if m.(model).quitting {
			fmt.Println("Installation canceled.")
			os.Exit(0)
		}

		// Extract the user input
		inputs := m.(model).inputs
		if len(inputs) < 2 {
			fmt.Println("Error: Not enough input provided.")
			os.Exit(1)
		}
		relayDomain := inputs[0]
		email := inputs[1]

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

