package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/azuki-bar/packtrack/output"
	"github.com/azuki-bar/packtrack/packagemanager"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

type config struct {
	Dryrun          bool
	WebHookEndpoint string
	Format          format
	Manager         packagemanager.Name
	// TODO: Colorは一部のアクターのみに使うコンフィグなのでトップレベルの設定にいるのはおかしいかも
	Color bool
}

func Main(in io.Reader, stdout, stderr io.Writer) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.packtrack/")
	viper.AddConfigPath("$XDG_CONFIG_PATH/.packtrack/")
	pflag.String("format", "default", "specify output format")
	pflag.Bool("dryrun", false, "(NOT IMPL) dryrun")
	pflag.String("manager", "", "specify package manager, from `brew, yay`")
	pflag.Bool("color", true, "output with color")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}
	c := config{}
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	packageManager := packagemanager.New(c.Manager, nil)
	outdated, err := packageManager.Outdated(ctx)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return err
	}
	actor := func() output.Actor {
		switch c.Format {
		case slackWebFookFormat:
			return output.NewSlack(output.SlackConf{Endpoint: c.WebHookEndpoint}, outdated)
		case jsonFormat:
			return output.NewJSON(output.JSONConfig{}, outdated)
		case plainFormat:
			fallthrough
		default:
			return output.NewPlain(output.PlainConfig{IsColor: c.Color}, outdated)
		}
	}()
	return actor.Exec(ctx, stdout, stderr)
}
