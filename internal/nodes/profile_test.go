package nodes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProfile(t *testing.T) {
	d := t.TempDir()
	valid := `[Interface]
PrivateKey = abc
Address = 10.8.0.2/32

[Peer]
PublicKey = def
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 203.0.113.10:51820
`
	path := filepath.Join(d, "valid.conf")
	if err := os.WriteFile(path, []byte(valid), 0600); err != nil { t.Fatal(err) }
	if err := ValidateProfile(path); err != nil { t.Fatalf("valid profile rejected: %v", err) }

	bad := filepath.Join(d, "bad.conf")
	if err := os.WriteFile(bad, []byte(valid+"# YOUR_SERVER_PUBLIC_KEY"), 0600); err != nil { t.Fatal(err) }
	if err := ValidateProfile(bad); err == nil { t.Fatal("placeholder profile accepted") }
}
