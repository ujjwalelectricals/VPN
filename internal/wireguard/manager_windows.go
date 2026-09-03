//go:build windows

package wireguard

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ujjwalelectricals/VPN/internal/model"
)

func executable() string {
	candidates := []string{
		`C:\Program Files\WireGuard\wireguard.exe`,
		`C:\Program Files (x86)\WireGuard\wireguard.exe`,
		`wireguard.exe`,
	}
	for _, c := range candidates {
		if c == `wireguard.exe` { return c }
		if _, err := exec.LookPath(c); err == nil { return c }
	}
	return candidates[0]
}

func serviceName(node model.Node) string {
	id := strings.NewReplacer("-", "", "_", "", " ", "").Replace(node.ID)
	return "WireGuardTunnel$" + id
}

func InstallAndStart(repoRoot string, node model.Node) error {
	cfg := filepath.Join(repoRoot, node.Config)
	if err := exec.Command(executable(), "/installtunnelservice", cfg).Run(); err != nil {
		// If the service already exists, starting it is enough.
		if !serviceExists(serviceName(node)) { return fmt.Errorf("install tunnel: %w", err) }
	}
	if err := exec.Command("sc.exe", "start", serviceName(node)).Run(); err != nil {
		// A service that is already running is fine.
		if !strings.Contains(strings.ToLower(string(err.Error())), "already") { return fmt.Errorf("start tunnel: %w", err) }
	}
	return nil
}

func Stop(node model.Node) error {
	cmd := exec.Command("sc.exe", "stop", serviceName(node))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stop tunnel: %w", err)
	}
	return nil
}

func Uninstall(node model.Node) error {
	_ = exec.Command("sc.exe", "stop", serviceName(node)).Run()
	if err := exec.Command(executable(), "/uninstalltunnelservice", strings.TrimPrefix(serviceName(node), "WireGuardTunnel$")).Run(); err != nil {
		return fmt.Errorf("uninstall tunnel: %w", err)
	}
	return nil
}

func serviceExists(name string) bool {
	return exec.Command("sc.exe", "query", name).Run() == nil
}

func Status(node model.Node) model.Status {
	if _, err := exec.LookPath("sc.exe"); err != nil {
		return model.Status{Installed: false, Message: "Windows service tooling is unavailable"}
	}
	if !serviceExists(serviceName(node)) {
		return model.Status{Installed: false, Tunnel: serviceName(node), Message: "Tunnel not installed"}
	}
	out, err := exec.Command("sc.exe", "query", serviceName(node)).CombinedOutput()
	if err != nil {
		return model.Status{Installed: true, Tunnel: serviceName(node), Message: string(out)}
	}
	connected := strings.Contains(strings.ToUpper(string(out)), "RUNNING")
	msg := "Stopped"
	if connected { msg = "Connected" }
	return model.Status{Installed: true, Connected: connected, Tunnel: serviceName(node), Message: msg}
}
