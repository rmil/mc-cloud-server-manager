package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var ErrKeyFileDeclined = errors.New("user declined key-file generation")

// KeyGenerationWizard prompts a user if they want to generate a key
// and generates a pair.
func KeyGenerationWizard(defaultFilePath string) (string, error) {
	filepath, err := askUserForKeyGeneration(defaultFilePath)
	if err != nil {
		return "", ErrKeyFileDeclined
	}

	err = generateKeyPair(filepath)
	if err != nil {
		return filepath, fmt.Errorf("failed to generate key-pair: %w", err)
	}
	return filepath, nil
}

// askUserForKeyGeneration will prompt the user if they want to generate
// an SSH key-pair and where they would like to save it.
func askUserForKeyGeneration(defaultFilePath string) (string, error) {
	filePath := ""
	fmt.Print("Would you like to generate a pair? [Y/n] ")
	if isYesResponse() {
		fmt.Printf("Enter file in which to save the key (%s): ", defaultFilePath)
		fmt.Scanln(&filePath)
		if filePath == "" {
			filePath = defaultFilePath
		}
	} else {
		return "", ErrKeyFileDeclined
	}
	return filePath, nil
}

// isYesResponse will check the user input if it is yes ("Y" / "y" / "")
// or no (anything else).
func isYesResponse() (ok bool) {
	response := ""
	fmt.Scanln(&response)
	if response == "Y" || response == "y" || response == "" {
		return true
	}
	return false
}

// generateKeyPair generates an RSA 4096 bit key-pair saving it at the
// filename name provided and a public key at the filename + ".pub".
func generateKeyPair(filename string) error {
	// Generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}
	publickey := &privatekey.PublicKey

	// Dump private key to file.
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		return fmt.Errorf("failed to encode private key file: %w", err)
	}

	// Dump public key to file.
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return fmt.Errorf("failed to dump public key: %w", err)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyFilename := fmt.Sprintf("%s.pub", filename)
	publicPem, err := os.Create(publicKeyFilename)
	if err != nil {
		return fmt.Errorf("failed to create public key file \"%s\": %w", publicKeyFilename, err)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		return fmt.Errorf("failed to encode public key file: %w", err)
	}
	return nil
}
