package output_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/azuki-bar/packtrack/output"
	"github.com/azuki-bar/packtrack/packagemanager"
	"github.com/stretchr/testify/assert"
)

func TestPlainExec(t *testing.T) {
	testCases := map[string]struct {
		appList []packagemanager.AppPackage
		stdout  string
		stdErr  string
		wantErr bool
	}{
		"no applist": {
			appList: []packagemanager.AppPackage{},
			stdout:  "not need to update!\n",
			stdErr:  "",
			wantErr: false,
		},
		"no applist by nil": {
			appList: nil,
			stdout:  "not need to update!\n",
			stdErr:  "",
			wantErr: false,
		},
		"exist one app in applist": {
			appList: []packagemanager.AppPackage{
				{Name: "packageA", LocalVersion: "1.1.1", RemoteVersion: "1.1.2"},
			},
			stdout: `Outdated Packages list
"packageA": 1.1.1 -> 1.1.2
`,
			stdErr:  "",
			wantErr: false,
		},
	}
	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {

			plain := output.NewPlain(output.PlainConfig{}, v.appList)
			stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")

			err := plain.Exec(context.Background(), stdout, stderr)
			assert.Equal(t, v.stdout, stdout.String())
			assert.Equal(t, v.stdErr, stderr.String())
			if v.wantErr {
				assert.Error(t, err)
			}
			assert.NoError(t, err)
		})
	}
}
