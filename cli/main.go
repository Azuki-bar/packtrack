package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/azuki-bar/packtrack/output"
	"github.com/azuki-bar/packtrack/packagemanager"
	"github.com/spf13/viper"
)

type format string

const (
	slackWebFookFormat format = "slack"
	jsonFormat         format = "json"
	plaintFormat       format = "plain"
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
		*f = plaintFormat
	}
	return nil
}

type config struct {
	Dryrun          bool
	WebHookEndpoint string
	Format          format
	Manager         string
}

func Main(in io.Reader, stdout, stderr io.Writer) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.packtrack/")
	viper.AddConfigPath("$XDG_CONFIG_PATH/.packtrack/")
	c := config{}
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	packageManager := packagemanager.New(c.Manager, nil)
	outdated, err := packageManager.Outdated(ctx)
	if err != nil {
		stderr.Write([]byte(err.Error()))
		return err
	}
	actor := func() output.Actor {
		switch c.Format {
		case slackWebFookFormat:
			return output.NewSlack(output.SlackConf{Endpoint: c.WebHookEndpoint}, outdated)
		case plaintFormat:
			fallthrough
		default:
			return output.NewPlain(output.PlainConfig{}, outdated)
		}
	}()
	return actor.Exec(ctx, stdout, stderr)
}