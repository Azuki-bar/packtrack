package packagemanager

import (
	"context"
	"fmt"
)

type AppPackage struct {
	Name          string `json:"name"`
	LocalVersion  string `json:"local_version"`
	RemoteVersion string `json:"remote_version"`
}

func (ap AppPackage) String() string {
	return fmt.Sprintf(`"%s": %s -> %s`, ap.Name, ap.LocalVersion, ap.RemoteVersion)
}

type packageManagerName string

const (
	brew           packageManagerName = "brew"
	yayPackageName packageManagerName = "yay"
)

type Manager struct {
	name    packageManagerName
	manager interface {
		Outdated(ctx context.Context) ([]AppPackage, error)
	}
}

func New(packageManager string, extraArgs []string) Manager {
	switch packageManager {
	case "brew":
		return Manager{name: brew, manager: homeBrew{extraArgs: extraArgs}}
	case "yay":
		return Manager{name: yayPackageName, manager: yay{}}
	default:
		return Manager{}
	}
}

func (m Manager) Outdated(ctx context.Context) ([]AppPackage, error) {
	if m.manager != nil {
		return m.manager.Outdated(ctx)
	}
	return nil, fmt.Errorf("unimplemented")
}
