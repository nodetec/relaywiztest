package network

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Function to print info messages
func printInfo(message string) {
	fmt.Println(message)
}

// Function to extract the directory name from the domain
func getDirectoryName(domainName string) string {
	domainParts := strings.Split(domainName, ".")
	if len(domainParts) > 2 {
		return domainParts[1]
	}
	return domainParts[0]
}

// Function to configure nginx for HTTP
func ConfigureNginxHttp(domainName string) {
	dirName := getDirectoryName(domainName)

	printInfo("Creating necessary directories...")
	err := os.MkdirAll(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", dirName), 0755)
	if err != nil {
		fmt.Printf("Error creating directories: %v\n", err)
		return
	}

	printInfo("Removing existing nginx configuration if it exists...")
	err = os.Remove("/etc/nginx/conf.d/nostr_relay.conf")
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error removing existing nginx configuration: %v\n", err)
		return
	}

	printInfo("Configuring nginx for HTTP...")
	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket {
    server 127.0.0.1:8080;
}

# %s
server {
    listen 80;
    listen [::]:80;
    server_name %s;

    location /.well-known/acme-challenge/ {
        root /var/www/%s;
        allow all;
    }

    location / {
        proxy_pass http://websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $remote_addr;
    }
}
`, domainName, domainName, dirName)

	err = os.WriteFile("/etc/nginx/conf.d/nostr_relay.conf", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("Error writing nginx configuration: %v\n", err)
		return
	}

	printInfo("Reloading nginx to apply the configuration...")
	err = exec.Command("systemctl", "reload", "nginx").Run()
	if err != nil {
		fmt.Printf("Error reloading nginx: %v\n", err)
		return
	}

	printInfo("Nginx configuration for HTTP completed successfully.")
}

