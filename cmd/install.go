package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nodetec/relaywiz/pkg/manager"
	"github.com/nodetec/relaywiz/pkg/network"
	"github.com/nodetec/relaywiz/pkg/relay"
	"github.com/spf13/cobra"
)

type model struct {
	textInput textinput.Model
	cursor    int
	inputs    []string
	quitting  bool
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Relay Domain"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	return model{
		textInput: ti,
		inputs:    make([]string, 0),
		quitting:  false,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.inputs = append(m.inputs, m.textInput.Value())
			m.textInput.Reset()
			if m.cursor == 0 {
				m.textInput.Placeholder = "Email Address"
				m.cursor++
			} else {
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}

	// case errMsg:
	// 	m.err = msg
	// 	return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Exiting the program. Goodbye!\n"
	}

	return fmt.Sprintf(
		"%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())

		// Run the Bubble Tea program to get user input
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if m.(model).quitting {
			fmt.Println("Installation canceled.")
			os.Exit(0)
		}

		// Extract the user inputs
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

