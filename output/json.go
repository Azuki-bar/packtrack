package output

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/azuki-bar/packtrack/packagemanager"
)

type JSON struct {
	appList []packagemanager.AppPackage
	config  JSONConfig
}
type JSONConfig struct {
	Indent string
}

func NewJSON(config JSONConfig, appList []packagemanager.AppPackage) *JSON {
	return &JSON{appList: appList, config: config}
}

func (j *JSON) Exec(ctx context.Context, stdout, stderr io.Writer) error {
	b, err := func() ([]byte, error) {
		if j.config.Indent == "" {
			return json.Marshal(j.appList)
		}
		return json.MarshalIndent(j.appList, "", j.config.Indent)
	}()
	if err != nil {
		return fmt.Errorf("marshaling JSON err, err=%w", err)
	}
	_, err = stdout.Write(b)
	return err
}
