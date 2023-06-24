package output

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/azuki-bar/packtrack/packagemanager"
	"github.com/samber/lo"
)

type Slack struct {
	conf     SlackConf
	outdated []packagemanager.AppPackage
}

type SlackConf struct {
	configBase
	Endpoint string
}

func NewSlack(conf Config, outdated []packagemanager.AppPackage) (*Slack, error) {
	c, ok := conf.(*SlackConf)
	if !ok {
		return nil, fmt.Errorf("type error")
	}
	return &Slack{conf: *c, outdated: outdated}, nil
}
func (s *Slack) Exec(ctx context.Context, stdout, stderr io.Writer) error {
	body, err := s.toJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.conf.Endpoint, bytes.NewReader(body))
	if err != nil {
		fmt.Fprint(stderr, "create http client failed")
		fmt.Fprint(stderr, err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		stderr.Write([]byte("execute post failed"))
		stderr.Write([]byte(err.Error()))
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		return err
	}
	_, err = stdout.Write([]byte("succeed to post"))
	return err
}

func (s *Slack) toJSON() ([]byte, error) {
	outdateStr := strings.Join(lo.Map(s.outdated, func(p packagemanager.AppPackage, _ int) string { return p.String() }), "\n")
	payload := struct {
		Text string `json:"text"`
	}{
		Text: outdateStr,
	}
	return json.Marshal(&payload)
}
