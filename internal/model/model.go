package model

type Node struct {
	ID       string `json:"id"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Name     string `json:"name"`
	Config   string `json:"config"`
	Endpoint string `json:"endpoint"`
	Enabled  bool   `json:"enabled"`
}

type Status struct {
	Installed bool   `json:"installed"`
	Connected bool   `json:"connected"`
	Tunnel    string `json:"tunnel"`
	Message   string `json:"message"`
}
