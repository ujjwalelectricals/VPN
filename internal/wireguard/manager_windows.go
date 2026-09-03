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
		if c == `wireguard.exe` {
			return c
		}
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}
	return candidates[0]
}

func serviceName(node model.Node) string {
	id := strings.NewReplacer("-", "", "_", "", " ", "").Replace(node.ID)
	return "WireGuardTunnel$" + id
}

func InstallAndStart(repoRoot string, node model.Node) error {
	cfg := filepath.Join(repoRoot, node.Config)
	if !serviceExists(serviceName(node)) {
		if err := exec.Command(executable(), "/installtunnelservice", cfg).Run(); err != nil {
			return fmt.Errorf("install tunnel: %w", err)
		}
	}
	if status := Status(node); status.Connected {
		return nil
	}
	out, err := exec.Command("sc.exe", "start", serviceName(node)).CombinedOutput()
	if err != nil {
		if status := Status(node); status.Connected {
			return nil
		}
		return fmt.Errorf("start tunnel: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func Stop(node model.Node) error {
	out, err := exec.Command("sc.exe", "stop", serviceName(node)).CombinedOutput()
	if err != nil && !strings.Contains(strings.ToUpper(string(out)), "STOPPED") {
		return fmt.Errorf("stop tunnel: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func Uninstall(node model.Node) error {
	_ = Stop(node)
	name := strings.TrimPrefix(serviceName(node), "WireGuardTunnel$")
	out, err := exec.Command(executable(), "/uninstalltunnelservice", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("uninstall tunnel: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func serviceExists(name string) bool {
	return exec.Command("sc.exe", "query", name).Run() == nil
}

func Status(node model.Node) model.Status {
	if _, err := exec.LookPath("sc.exe"); err != nil {
		return model.Status{Installed: false, Message: "Windows service tooling unavailable"}
	}
	if !serviceExists(serviceName(node)) {
		return model.Status{Installed: false, Tunnel: serviceName(node), Message: "Tunnel not installed"}
	}
	out, err := exec.Command("sc.exe", "query", serviceName(node)).CombinedOutput()
	if err != nil {
		return model.Status{Installed: true, Tunnel: serviceName(node), Message: string(out)}
	}
	connected := strings.Contains(strings.ToUpper(string(out)), "RUNNING")
	message := "Stopped"
	if connected {
		message = "Connected"
	}
	return model.Status{Installed: true, Connected: connected, Tunnel: serviceName(node), Message: message}
}
