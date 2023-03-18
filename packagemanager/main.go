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

type Name string

const (
	brew           Name = "brew"
	yayPackageName Name = "yay"
)

func (n *Name) UnmarshalMap(input any) error {
	v, ok := input.(string)
	if !ok {
		return fmt.Errorf("no supported package manager")
	}
	switch v {
	case "homebrew", "brew":
		*n = brew
	case "yay":
		*n = yayPackageName
	default:
	}
	return nil
}

type Manager struct {
	name    Name
	manager interface {
		Outdated(ctx context.Context) ([]AppPackage, error)
	}
}

func New(packageManager Name, extraArgs []string) Manager {
	switch packageManager {
	case brew:
		return Manager{name: brew, manager: homeBrew{extraArgs: extraArgs}}
	case yayPackageName:
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
