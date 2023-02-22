package packagemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/samber/lo"
)

type brewPackage struct {
	PackageName      string `json:"name"`
	InstalledVersion string `json:"installed_version"`
	Version          string `json:"current_version"`
}

type homeBrew struct {
	extraArgs []string
}

func (h homeBrew) Outdated(ctx context.Context) ([]AppPackage, error) {
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
	type brewPackages []brewPackage
	ps := new(brewPackages)
	if err := json.Unmarshal(b, ps); err != nil {
		return nil, fmt.Errorf("unmarshal json from brew failed, err=`%w`", err)
	}
	return lo.Map([]brewPackage(*ps), func(pac brewPackage, _ int) AppPackage {
		return AppPackage{
			Name:          pac.PackageName,
			LocalVersion:  pac.InstalledVersion,
			RemoteVersion: pac.Version,
		}
	}), nil
}
