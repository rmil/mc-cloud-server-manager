package cloud

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

var ErrDCNotFound = errors.New("dc not found")

func (c *Clouder) ConfigureCloud(ctx context.Context, publicKey ssh.PublicKey) error {
	err := c.selectDatacentre()
	if err != nil {
		return fmt.Errorf("failed to select datacentre: %w", err)
	}

	keyFingerprint := getKeyFingerprint(publicKey)

	err = c.checkForKey(ctx, keyFingerprint)
	if err != nil {
		if errors.As(err, &ErrNoFingerprintsMatch) {
			fmt.Println("Public key not present on Hetzner")
			err = c.uploadPublicKey(ctx, publicKey)
			if err != nil {
				return fmt.Errorf("failed to upload private key to Hetzner: %w", err)
			}
			fmt.Printf("Successfully uploaded key \"%s\"\n", keyFingerprint)
		} else {
			return fmt.Errorf("failed to check for key: %w", err)
		}
	} else {
		fmt.Println("SSH public-key present on Hetzner")
	}

	return nil
}

func (c *Clouder) selectDatacentre() error {
	manualSelection := false
	dc, err := c.getRecommendedDatacentre()
	if err != nil {
		if errors.As(err, &ErrRecommendationNotFound) {
			fmt.Println("Recommendation not found, manual selection required")
			manualSelection = true
		}
	}
	selectedDC := dc.Name
	fmt.Printf(`Select the recommended datacentre? (%s | %s) [Y/n] `, dc.Name, dc.Location.Description)

	if !manualSelection && isYesResponse() {
		c.conf.selectedDatacentre = dc.ID
	} else {
		for {
			fmt.Println("Manual datacentre selection:")

			for _, dc := range c.dc.Datacentres {
				fmt.Printf("%s\n", dc.Name)
			}

			fmt.Print("\nPlease enter a DC name: ")
			dcName := ""
			fmt.Scanln(&dcName)

			dc, err := c.getDCByName(dcName)
			if err == nil {
				c.conf.selectedDatacentre = dc.ID
				selectedDC = dc.Name
				break
			}
		}
	}

	fmt.Printf("Selected %s\n", selectedDC)

	return nil
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
