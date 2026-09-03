package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ujjwalelectricals/VPN/internal/nodes"
	"github.com/ujjwalelectricals/VPN/internal/wireguard"
)

func repoRoot() string {
	exe, err := os.Executable()
	if err == nil {
		if p := filepath.Clean(filepath.Join(filepath.Dir(exe), "..", "..")); fileExists(filepath.Join(p, "nodes.json")) { return p }
	}
	cwd, _ := os.Getwd()
	return cwd
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }

func main() {
	root := repoRoot()
	if len(os.Args) < 2 { usage(); return }
	switch os.Args[1] {
	case "list": cmdList(root)
	case "connect": requireArg(os.Args, 2, "vpnctl connect <node-id>"); cmdConnect(root, os.Args[2])
	case "disconnect": requireArg(os.Args, 2, "vpnctl disconnect <node-id>"); cmdDisconnect(root, os.Args[2])
	case "status": requireArg(os.Args, 2, "vpnctl status <node-id>"); cmdStatus(root, os.Args[2])
	case "doctor": cmdDoctor(root)
	case "dashboard": cmdDashboard(root)
	case "version": fmt.Println("Ujjwal FreeVPN 0.1.0")
	default: usage()
	}
}

func requireArg(args []string, index int, usage string) { if len(args) <= index { fatal(usage) } }

func cmdList(root string) {
	ns, err := nodes.Load(root); if err != nil { fatal(err.Error()) }
	for _, n := range ns {
		state := "disabled"; if n.Enabled { state = "ready" }
		fmt.Printf("%-14s %-16s %-18s %-8s %s\n", n.ID, n.Country, n.City, state, n.Name)
	}
}

func cmdConnect(root, id string) {
	n, err := nodes.Find(root, id); if err != nil { fatal(err.Error()) }
	if !n.Enabled { fatal("node is disabled; add a valid WireGuard configuration before connecting") }
	if err := wireguard.InstallAndStart(root, n); err != nil { fatal(err.Error()) }
	fmt.Printf("Connected to %s (%s, %s).\n", n.Name, n.City, n.Country)
}

func cmdDisconnect(root, id string) {
	n, err := nodes.Find(root, id); if err != nil { fatal(err.Error()) }
	if err := wireguard.Stop(n); err != nil { fatal(err.Error()) }
	fmt.Printf("Disconnected %s.\n", n.Name)
}

func cmdStatus(root, id string) {
	n, err := nodes.Find(root, id); if err != nil { fatal(err.Error()) }
	printJSON(wireguard.Status(n))
}

func cmdDoctor(root string) {
	report := map[string]any{
		"os": runtime.GOOS,
		"arch": runtime.GOARCH,
		"repo_root": root,
		"nodes_catalog": fileExists(filepath.Join(root, "nodes.json")),
		"profiles_dir": fileExists(filepath.Join(root, "profiles")),
		"wireguard_manager": runtime.GOOS == "windows",
		"note": "Add a valid WireGuard client profile for each enabled node. Never commit real private keys.",
	}
	printJSON(report)
}

func cmdDashboard(root string) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepath.Join(root, "web"))))
	mux.HandleFunc("/api/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		ns, err := nodes.Load(root); if err != nil { http.Error(w, err.Error(), 500); return }
		writeJSON(w, ns)
	})
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		writeJSON(w, map[string]any{"ok":true,"os":runtime.GOOS})
	})
	mux.HandleFunc("/api/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		id := strings.TrimSpace(r.URL.Query().Get("id")); if id == "" { http.Error(w, "missing id", 400); return }
		n, err := nodes.Find(root, id); if err != nil { http.Error(w, err.Error(), 404); return }
		if !n.Enabled { http.Error(w, "node is disabled", 409); return }
		if err := wireguard.InstallAndStart(root, n); err != nil { http.Error(w, err.Error(), 500); return }
		writeJSON(w, map[string]any{"ok":true,"message":"connected","node":n.ID})
	})
	mux.HandleFunc("/api/disconnect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		id := strings.TrimSpace(r.URL.Query().Get("id")); if id == "" { http.Error(w, "missing id", 400); return }
		n, err := nodes.Find(root, id); if err != nil { http.Error(w, err.Error(), 404); return }
		if err := wireguard.Stop(n); err != nil { http.Error(w, err.Error(), 500); return }
		writeJSON(w, map[string]any{"ok":true,"message":"disconnected","node":n.ID})
	})
	addr := "127.0.0.1:7070"
	fmt.Printf("Dashboard: http://%s (local machine only)\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil { fatal(err.Error()) }
}

func usage() {
	fmt.Println(strings.TrimSpace(`vpnctl commands:
  list                    list configured locations
  connect <node-id>       install/start a WireGuard tunnel
  disconnect <node-id>    stop a WireGuard tunnel
  status <node-id>        show tunnel status
  doctor                  check local prerequisites
  dashboard               open the local control dashboard
  version                 print version

The project manages your own WireGuard nodes. It does not provide paid or third-party VPN service.`))
}

func printJSON(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
func writeJSON(w http.ResponseWriter, v any) { w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(v) }
func fatal(msg string) { fmt.Fprintln(os.Stderr, "error:", msg); os.Exit(1) }
