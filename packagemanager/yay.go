package packagemanager

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
)

type yay struct{}

var yayPackageRegexp = regexp.MustCompile(`(.*) (.+) -> (.+)`)

func (y yay) Outdated(ctx context.Context) ([]AppPackage, error) {
	cmd := exec.Command("yay", "-Qu")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	aps := make([]AppPackage, 0)
	r := bufio.NewScanner(bytes.NewReader(b))
	for r.Scan() {
		text := r.Text()
		parsed := yayPackageRegexp.FindStringSubmatch(text)
		if len(parsed) < 4 {
			return nil, fmt.Errorf("parse yay packages failed, str=%s", text)
		}
		aps = append(aps, AppPackage{Name: parsed[1], LocalVersion: parsed[2], RemoteVersion: parsed[3]})
	}
	return aps, nil
}
