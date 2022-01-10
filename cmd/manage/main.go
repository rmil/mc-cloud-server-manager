package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rmil/mc-cloud-server-manager/cloud"
	"github.com/rmil/mc-cloud-server-manager/config"
	"github.com/rmil/mc-cloud-server-manager/ssh"
)

func main() {
	// Load environment.
	filePath := flag.String("c", "config.json", "config file location")
	flag.Parse()

	// Load configuration.
	globalConfig, err := config.LoadConfig(*filePath)
	if err != nil {
		if errors.As(err, &os.ErrNotExist) {
			err = config.GenerateDefaultConfigFile(*filePath)
			if err != nil {
				log.Fatalf("Failed to make config file: %+v", err)
			}
			fmt.Printf("Generated config file \"%s\", please update the file with the correct information\n", *filePath)
			os.Exit(0)
		}
		log.Fatalf("Failed to load configuration file: %+v", err)
	}

	// Setup SSH Key-pair.
	_, err = ssh.GetKey(globalConfig.SSH.KeyFilePath)
	if err != nil {
		if errors.As(err, &ssh.ErrNoKeyFound) {
			fmt.Println("failed to find SSH key")
			generatedKeyName, err := ssh.KeyGenerationWizard(globalConfig.SSH.KeyFilePath)
			if err != nil {
				if errors.As(err, &ssh.ErrKeyFileDeclined) {
					fmt.Println("Generate a key-pair using `ssh-keygen` and try again")
					os.Exit(0)
				}
				log.Fatalf("key generation wizard failed: %+v", err)
			}
			globalConfig.SSH.KeyFilePath = generatedKeyName
			fmt.Printf("Successfully generated key \"%s\"\n", generatedKeyName)
		} else {
			log.Fatalf("failed to get SSH key: %+v", err)
		}
	}

	// Re-getting in case a new one was generated.
	key, err := ssh.GetKey(globalConfig.SSH.KeyFilePath)
	if err != nil {
		log.Fatalf("failed to get SSH key: %+v", err)
	}

	// Connect to Hetzner
	cloudClient, err := cloud.NewCloud(globalConfig.Cloud)
	if err != nil {
		log.Fatalf("failed to create Hetzner connection: %+v", err)
	}

	err = cloudClient.ConfigureCloud(context.Background(), key.PublicKey())
	if err != nil {
		log.Fatalf("failed to configure Hetzner: %+v", err)
	}

	// Connect to SSH server.
	sshClient, err := ssh.Connect(globalConfig.SSH)
	if err != nil {
		log.Fatalf("failed to connect to SSH server: %+v", err)
	}

	sshClient.GetServerInfo()
}
