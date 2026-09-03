//go:build !windows

package wireguard

import (
	"errors"
	"github.com/ujjwalelectricals/VPN/internal/model"
)

func InstallAndStart(_ string, _ model.Node) error { return errors.New("VPN tunnel management is currently supported on Windows only") }
func Stop(_ model.Node) error { return errors.New("VPN tunnel management is currently supported on Windows only") }
func Uninstall(_ model.Node) error { return errors.New("VPN tunnel management is currently supported on Windows only") }
func Status(_ model.Node) model.Status { return model.Status{Installed:false, Message:"Windows-only WireGuard manager"} }
