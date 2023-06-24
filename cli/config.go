package cli

import (
	"fmt"

	"github.com/azuki-bar/packtrack/output"
	"github.com/azuki-bar/packtrack/packagemanager"
)

type format string

const (
	slackWebFookFormat format = "slack"
	jsonFormat         format = "json"
	plainFormat        format = "plain"
	templateFormat     format = "template"
)

func (f *format) UnmarshalMap(input any) error {
	v, ok := input.(string)
	if !ok {
		return fmt.Errorf("no supported value type")
	}
	switch v {
	case "slack":
		*f = slackWebFookFormat
	case "json":
		*f = jsonFormat
	default:
		*f = plainFormat
	}
	return nil
}

type Config struct {
	Dryrun  bool
	Format  format
	Manager packagemanager.Name
	// TODO: Colorは一部のアクターのみに使うコンフィグなのでトップレベルの設定にいるのはおかしいかも
	OutputOption map[format]output.Config
}
