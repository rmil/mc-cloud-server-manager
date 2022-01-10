package cloud

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/crypto/ssh"
)

var ErrNoFingerprintsMatch = errors.New("no key fingerprints matched")

func (c *Clouder) checkForKey(ctx context.Context, clientKeyFingerprint string) error {
	keys, err := c.client.SSHKey.All(ctx)
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}
	for _, key := range keys {
		if key.Fingerprint == clientKeyFingerprint {
			return nil
		}
	}

	return ErrNoFingerprintsMatch
}

func getKeyFingerprint(key ssh.PublicKey) string {
	return ssh.FingerprintLegacyMD5(key)
}

func (c *Clouder) uploadPublicKey(ctx context.Context, publicKey ssh.PublicKey) error {
	formattedKey := fmt.Sprintf("ssh-rsa %s", base64.StdEncoding.EncodeToString(publicKey.Marshal()))

	_, _, err := c.client.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
		Name:      c.conf.AppName,
		PublicKey: formattedKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create public key: %w", err)
	}
	return nil
}
