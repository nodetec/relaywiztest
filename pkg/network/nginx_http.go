package network

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"github.com/nodetec/relaywiz/pkg/utils"
)

// Function to configure nginx for HTTP
func ConfigureNginxHttp(domainName string) {
	dirName := utils.GetDirectoryName(domainName)

	fmt.Println("Creating necessary directories...")
	err := os.MkdirAll(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", dirName), 0755)
	if err != nil {
		log.Fatalf("Error creating directories: %v", err)
	}

	fmt.Println("Removing existing nginx configuration if it exists...")
	err = os.Remove("/etc/nginx/conf.d/nostr_relay.conf")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing existing nginx configuration: %v", err)
	}

	fmt.Println("Configuring nginx for HTTP...")
	configContent := fmt.Sprintf(`map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

upstream websocket {
    server 0.0.0.0:3334;
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
		log.Fatalf("Error writing nginx configuration: %v", err)
	}

	fmt.Println("Reloading nginx to apply the configuration...")
	err = exec.Command("systemctl", "reload", "nginx").Run()
	if err != nil {
		log.Fatalf("Error reloading nginx: %v", err)
	}

	fmt.Println("Nginx configuration for HTTP completed successfully.")
}

