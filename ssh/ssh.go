package ssh

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type (
	Config struct {
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		KeyFilePath string `json:"keyfilePath"`
	}
	Client struct {
		ssh *ssh.Session
	}
)

var ErrNoKeyFound = errors.New("no private key found")

// Connect to an SSH server using a username and a SSH key.
func Connect(conf Config) (*Client, error) {
	signer, err := GetKey(conf.KeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", conf.Hostname, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Client{ssh: session}, nil
}

// GetKey will fetch the private key at the specified filepath.
func GetKey(filepath string) (ssh.Signer, error) {
	key, err := os.ReadFile(filepath)
	if err != nil {
		if errors.As(err, &os.ErrNotExist) {
			return nil, ErrNoKeyFound
		}
		return nil, fmt.Errorf("failed to read private-key file: %w", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private-key: %w", err)
	}

	return signer, nil
}
