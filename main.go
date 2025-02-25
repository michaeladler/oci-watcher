// SPDX-FileCopyrightText: 2025 Margo
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileContributor: Michael Adler <michael.adler@siemens.com>

package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/regclient/regclient"
	"golang.org/x/term"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	rc          *regclient.RegClient
)

func main() {
	configPath := path.Join(os.Getenv("HOME"), ".docker", "config.json")
	if !fileExists(configPath) {
		_ = os.MkdirAll(path.Dir(configPath), 0o755)

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Github username: ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Enter Github token (scope read:packages): ")
		passwordBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Print("\n")
		password := string(passwordBytes)

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		authConfig := fmt.Sprintf(`{
	"auths": {
		"ghcr.io": {
			"auth": "%s"
		}
	}
}`, encodedAuth)

		file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(authConfig)
		if err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}

	deployDir := flag.String("deployDir", "./deploy", "Directory to deploy")
	ociRegistry := flag.String("ociRegistry", "ghcr.io/silvanoc/poc-deploy:desired", "OCI registry URL")
	flag.Parse()

	rc = regclient.New(regclient.WithDockerCerts(), regclient.WithDockerCreds())

	defer cancel()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	running := true
	for running {
		select {
		case <-ticker.C:
			if err := reconcileDeployments(*ociRegistry, *deployDir); err != nil {
				log.Println("ERROR:", err)
			}
		case <-sigChan:
			log.Println("Exiting gracefully...")
			cancel()
			running = false
		}
	}
	log.Println("Bye")
}
