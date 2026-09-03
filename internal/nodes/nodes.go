package nodes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ujjwalelectricals/VPN/internal/model"
)

func Load(root string) ([]model.Node, error) {
	path := filepath.Join(root, "nodes.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var nodes []model.Node
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("decode nodes.json: %w", err)
	}
	return nodes, nil
}

func Find(root, id string) (model.Node, error) {
	ns, err := Load(root)
	if err != nil {
		return model.Node{}, err
	}
	for _, n := range ns {
		if n.ID == id {
			return n, nil
		}
	}
	return model.Node{}, fmt.Errorf("node %q not found", id)
}
