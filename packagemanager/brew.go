package packagemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type homeBrew struct {
	extraArgs []string
}

type brewPackages []brewPackage
type brewPackage struct {
	PackageName      string `json:"name"`
	InstalledVersion string `json:"installed_version"`
	Version          string `json:"current_version"`
}

func (bp brewPackage) Name() string          { return bp.PackageName }
func (bp brewPackage) LocalVersion() string  { return bp.InstalledVersion }
func (bp brewPackage) RemoteVersion() string { return bp.Version }

func (h homeBrew) Outdated(ctx context.Context) (*packages[brewPackage], error) {
	// #nosec G204
	args := []string{"outdated", "--json=v2"}
	args = append(args, h.extraArgs...)
	cmd := exec.CommandContext(ctx, "brew", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("exec brew command failed,err=`%w`", err)
	}
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("get exec output failed, err=`%w`", err)
	}
	ps := new(brewPackages)
	if err := json.Unmarshal(b, ps); err != nil {
		return nil, fmt.Errorf("unmarshal json from brew failed, err=`%w`", err)
	}
	return &packages[brewPackage]{packages: *ps}, nil
}
