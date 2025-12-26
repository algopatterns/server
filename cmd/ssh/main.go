package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/algorave/server/internal/ssh"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	cfg := ssh.DefaultConfig()

	if port := os.Getenv("ALGORAVE_SSH_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return fmt.Errorf("invalid SSH port: %w", err)
		}

		cfg.Port = p
	}

	if hostKey := os.Getenv("ALGORAVE_SSH_HOST_KEY"); hostKey != "" {
		cfg.HostKeyPath = hostKey
	}

	if maxConn := os.Getenv("ALGORAVE_SSH_MAX_CONNECTIONS"); maxConn != "" {
		mc, err := strconv.Atoi(maxConn)
		if err != nil {
			return fmt.Errorf("invalid max connections: %w", err)
		}

		cfg.MaxConnections = mc
	}

	cfg.ProductionMode = true

	server, err := ssh.NewServer(cfg)
	if err != nil {
		return err
	}

	return server.Start()
}
