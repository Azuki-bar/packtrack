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

func Main(in io.Reader, stdout, stderr io.Writer) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.packtrack/")
	viper.AddConfigPath("$XDG_CONFIG_PATH/.packtrack/")
	pflag.String("format", "default", "specify output format")
	pflag.Bool("dryrun", false, "(NOT IMPL) dryrun")
	pflag.String("manager", "", "specify package manager, from `brew, yay`")
	pflag.Bool("color", true, "output with color")
	pflag.String("indent", "", "json")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}
	c := Config{OutputOption: make(map[format]output.Config)}
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
	actor, err := func() (output.Actor, error) {
		switch c.Format {
		case slackWebFookFormat:
			return output.NewSlack(c.OutputOption[slackWebFookFormat], outdated)
		case jsonFormat:
			return output.NewJSON(c.OutputOption[jsonFormat], outdated)
		case plainFormat:
			fallthrough
		default:
			return output.NewPlain(c.OutputOption[plainFormat], outdated)
		}
	}()
	if err != nil {
		return fmt.Errorf("parse format failed, err=%w", err)
	}
	return actor.Exec(ctx, stdout, stderr)
}
