package output

import (
	"context"
	"fmt"
	"io"

	"github.com/azuki-bar/packtrack/packagemanager"
)

type Plain struct {
	appList []packagemanager.AppPackage
	config  PlainConfig
}
type PlainConfig struct {
	isColor bool
}

func NewPlain(config PlainConfig, appList []packagemanager.AppPackage) *Plain {
	return &Plain{
		appList: appList,
		config:  config,
	}
}

func (p *Plain) Exec(ctx context.Context, stdout, stderr io.Writer) error {
	if p.appList == nil || len(p.appList) == 0 {
		fmt.Fprintln(stdout, "not need to update!")
		return nil
	}
	stdout.Write([]byte("Outdated Packages list\n"))
	for _, v := range p.appList {
		_, err := stdout.Write([]byte(v.String() + "\n"))
		if err != nil {
			stderr.Write([]byte("something error occured in printing method\n"))
			stderr.Write([]byte(err.Error()))
			return err
		}
	}
	return nil
}
