package packagemanager

import (
	"context"
	"fmt"
)

type packageManagerName string
type AppPackage struct {
	Name          string `json:"name"`
	LocalVersion  string `json:"local_version"`
	RemoteVersion string `json:"remote_version"`
}

func (ap AppPackage) String() string {
	return fmt.Sprintf(`"%s": %s -> %s`, ap.Name, ap.LocalVersion, ap.RemoteVersion)
}

const (
	brew packageManagerName = "brew"
	yay  packageManagerName = "yay"
)

type Manager struct {
	name packageManagerName
	brew homeBrew
	yay  string
}

func (m Manager) Outdated(ctx context.Context) ([]AppPackage, error) {
	switch m.name {
	case brew:
		ps, err := m.brew.Outdated(ctx)
		if err != nil {
			return nil, err
		}
		return mapPackages(ps.packages), err
	case yay:
		fallthrough
	default:
		return nil, fmt.Errorf("unimplemented")
	}
}

func mapPackages[T _package](packages []T) []AppPackage {
	aps := make([]AppPackage, len(packages))
	for i, v := range packages {
		aps[i] = AppPackage{
			Name:          v.Name(),
			LocalVersion:  v.LocalVersion(),
			RemoteVersion: v.RemoteVersion(),
		}
	}
	return aps
}

type packages[T _package] struct {
	packages []T
}
type _package interface {
	Name() string
	LocalVersion() string
	RemoteVersion() string
}

func New(packageManager string, extraArgs []string) Manager {
	switch packageManager {
	case "brew":
		return Manager{name: brew, brew: homeBrew{extraArgs: extraArgs}}
	case "yay":
		return Manager{name: yay, yay: "not impl"}
	default:
		return Manager{}
	}
}
