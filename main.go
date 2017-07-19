package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/chrisurwin/ds-remove/agent"
	"github.com/urfave/cli"
)

var VERSION = "v0.1.0-dev"

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "ds-remove"
	app.Version = VERSION
	app.Usage = "Removes an AWS host (if in an ASG) when the disk space is nearly full"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug,d",
			Usage:  "Debug logging",
			EnvVar: "DEBUG",
		},
		/*		cli.StringFlag{
					Name:   "aws-key,k",
					Usage:  "AWS Access Key",
					EnvVar: "AWS_ACCESS_KEY_ID",
				},
				cli.StringFlag{
					Name:   "aws-secret,s",
					Usage:  "AWS Secret Key",
					EnvVar: "AWS_SECRET_ACCESS_KEY",
				},*/
		cli.DurationFlag{
			Name:   "poll-interval,i",
			Value:  5 * time.Second,
			Usage:  "Polling interval for checks",
			EnvVar: "POLL_INTERVAL",
		},
		cli.StringFlag{
			Name:   "arn,a",
			Usage:  "AWS Role ARN",
			EnvVar: "ARN",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) {
	/*if c.String("aws-key") == "" {
		log.Fatal("AWS Access Key required")
	}
	if c.String("aws-secret") == "" {
		log.Fatal("AWS Secret Key required")
	}*/
	a := agent.NewAgent(c.Duration("poll-interval"), c.String("arn"))
	a.Start()
}
