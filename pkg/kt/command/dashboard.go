package command

import (
	"github.com/rs/zerolog/log"
	"github.com/skratchdot/open-golang/open"
	"github.com/alibaba/kt-connect/pkg/kt/exec"
	"github.com/alibaba/kt-connect/pkg/kt/exec/kubectl"
	"github.com/alibaba/kt-connect/pkg/kt/options"
	"github.com/rs/zerolog"
	"github.com/urfave/cli"
)

// newDashboardCommand dashboard command
func newDashboardCommand(options *options.DaemonOptions) cli.Command {
	return cli.Command{
		Name:  "dashboard",
		Usage: "kt-connect dashboard",
		Subcommands: []cli.Command{
			{
				Name:  "init",
				Usage: "install/update dashboard to cluster",
				Action: func(c *cli.Context) error {
					if options.Debug {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}
					action := Action{}
					return action.ApplyDashboard()
				},
			},
			{
				Name:  "open",
				Usage: "open dashboard",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:        "port,p",
						Value:       "8080",
						Usage:       "port-forward kt dashboard to port",
						Destination: &options.DashboardOptions.Port,
					},
				},
				Action: func(c *cli.Context) error {
					if options.Debug {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}
					action := Action{}
					return action.OpenDashboard(options)
				},
			},
		},
	}
}


func (action *Action) ApplyDashboard() (err error) {
	command := kubectl.ApplyDashboardToCluster()
	log.Info().Msg("Install/Upgrade Dashboard to cluster")
	err = exec.RunAndWait(command, "apply kt dashboard", true)
	if err != nil {
		log.Error().Msg("Fail to apply dashboard, please check the log")
		return
	}
	return
}

func (action *Action) OpenDashboard(options *options.DaemonOptions) (err error) {
	ch := SetUpWaitingChannel()
	command := kubectl.PortForwardDashboardToLocal(options.DashboardOptions.Port)
	err = exec.BackgroundRun(command, "forward dashboard to localhost", true)
	if err != nil {
		return
	}
	err = open.Run("http://127.0.0.1:" + options.DashboardOptions.Port)

	s := <-ch
	log.Info().Msgf("Terminal Signal is %s", s)
	return
}