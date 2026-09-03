package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ujjwalelectricals/VPN/internal/model"
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
	if runtime.GOOS != "windows" {
		fmt.Println("VPN client control is Windows-first; node catalog and diagnostics remain portable.")
	}
	root := repoRoot()
	if len(os.Args) < 2 { usage(); return }
	switch os.Args[1] {
	case "list": cmdList(root)
	case "connect":
		if len(os.Args) != 3 { fatal("usage: vpnctl connect <node-id>") }
		cmdConnect(root, os.Args[2])
	case "disconnect":
		if len(os.Args) != 3 { fatal("usage: vpnctl disconnect <node-id>") }
		cmdDisconnect(root, os.Args[2])
	case "status":
		if len(os.Args) != 3 { fatal("usage: vpnctl status <node-id>") }
		cmdStatus(root, os.Args[2])
	case "doctor": cmdDoctor(root)
	case "dashboard": cmdDashboard(root)
	case "version": fmt.Println("Ujjwal FreeVPN 0.1.0")
	default: usage()
	}
}

func cmdList(root string) {
	ns, err := nodes.Load(root); if err != nil { fatal(err.Error()) }
	for _, n := range ns {
		state := "disabled"; if n.Enabled { state = "ready" }
		fmt.Printf("%-14s %-12s %-18s %-8s %s\n", n.ID, n.Country, n.City, state, n.Name)
	}
}

func cmdConnect(root, id string) {
	n, err := nodes.Find(root, id); if err != nil { fatal(err.Error()) }
	if !n.Enabled { fatal("node is disabled; add your own server details before connecting") }
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
		"note": "A node cannot be used until you provide a valid WireGuard client configuration and a reachable server.",
	}
	printJSON(report)
}

func cmdDashboard(root string) {
	addr := "127.0.0.1:7070"
	fs := http.FileServer(http.Dir(filepath.Join(root, "web")))
	http.Handle("/", fs)
	http.HandleFunc("/api/nodes", func(w http.ResponseWriter, _ *http.Request) {
		ns, err := nodes.Load(root); if err != nil { http.Error(w, err.Error(), 500); return }
		writeJSON(w, ns)
	})
	http.HandleFunc("/api/health", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, map[string]any{"ok":true,"os":runtime.GOOS}) })
	fmt.Printf("Dashboard: http://%s (local machine only)\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil { fatal(err.Error()) }
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

Important: this project supplies the client/control layer. You must provide your own legal WireGuard server endpoint/configuration.`))
}

func printJSON(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
func writeJSON(w http.ResponseWriter, v any) { w.Header().Set("Content-Type", "application/json"); json.NewEncoder(w).Encode(v) }
func fatal(msg string) { fmt.Fprintln(os.Stderr, "error:", msg); os.Exit(1) }
var _ = flag.CommandLine
var _ = model.Node{}
