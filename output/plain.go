package output

import (
	"context"
	"fmt"
	"io"

	"github.com/azuki-bar/packtrack/packagemanager"
	"github.com/fatih/color"
)

type Plain struct {
	appList []packagemanager.AppPackage
	config  PlainConfig
}
type PlainConfig struct {
	configBase
	IsColor bool
}

func NewPlain(config Config, appList []packagemanager.AppPackage) (*Plain, error) {
	c, ok := config.(*PlainConfig)
	if !ok {
		return nil, fmt.Errorf("type error")
	}
	return &Plain{appList: appList, config: *c}, nil
}

func (p *Plain) Exec(ctx context.Context, stdout, stderr io.Writer) error {
	if p.appList == nil || len(p.appList) == 0 {
		fmt.Fprintln(stdout, "not need to update!")
		return nil
	}
	color.NoColor = !p.config.IsColor

	fmt.Fprintln(stdout, "Outdated Packages list")
	fmt.Fprintln(stdout, "======")

	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgHiRed)
	for _, v := range p.appList {
		bold.Fprint(stdout, v.Name)
		fmt.Fprint(stdout, "\t")
		red.Fprint(stdout, v.LocalVersion)
		fmt.Fprint(stdout, " -> ")
		green.Fprint(stdout, v.RemoteVersion)
		fmt.Fprint(stdout, "\n")
	}
	return nil
}
