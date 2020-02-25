package main

import (
	"os"
	"github.com/alibaba/kt-connect/pkg/kt/command"
	opt "github.com/alibaba/kt-connect/pkg/kt/options"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	options := opt.NewDaemonOptions()

	app := cli.NewApp()
	app.Name = "KT Connect"
	app.Usage = ""
	app.Version = "master"
	app.Authors = command.NewCliAuthor()
	app.Flags = command.AppFlags(options)
	app.Commands = command.NewCommands(options)

	err := app.Run(os.Args)
	if err != nil {
		log.Error().Msg(err.Error())
		command.CleanupWorkspace(options)
		os.Exit(-1)
	}

}
