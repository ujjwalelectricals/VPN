package nodes

import (
	"fmt"
	"os"
	"strings"
)

func ValidateProfile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(data)
	checks := []string{"[Interface]", "PrivateKey = ", "[Peer]", "PublicKey = ", "Endpoint = "}
	for _, check := range checks {
		if !strings.Contains(text, check) {
			return fmt.Errorf("WireGuard profile is missing %q", check)
		}
	}
	if strings.Contains(text, "YOUR_") || strings.Contains(text, "REPLACE_ME") {
		return fmt.Errorf("WireGuard profile still contains placeholders")
	}
	return nil
}
