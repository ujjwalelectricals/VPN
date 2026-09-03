package nodes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndFind(t *testing.T) {
	d := t.TempDir()
	data := `[{"id":"india","country":"India","city":"Test","name":"India Exit","config":"profiles/india.conf","endpoint":"127.0.0.1:51820","enabled":false}]`
	if err := os.WriteFile(filepath.Join(d, "nodes.json"), []byte(data), 0644); err != nil { t.Fatal(err) }
	ns, err := Load(d); if err != nil { t.Fatal(err) }
	if len(ns) != 1 || ns[0].ID != "india" { t.Fatalf("unexpected nodes: %+v", ns) }
	n, err := Find(d, "india"); if err != nil || n.Country != "India" { t.Fatalf("find failed: %+v %v", n, err) }
	if _, err := Find(d, "missing"); err == nil { t.Fatal("expected missing node error") }
}
