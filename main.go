package main

import (
	"log"
	"os"

	"github.com/alex-way/changesets/cmd/add"
	"github.com/alex-way/changesets/cmd/version"
	"github.com/urfave/cli/v2"
)

var addFlags = []cli.Flag{
	&cli.StringFlag{Name: "bump-type", Aliases: []string{"t"}},
	&cli.StringFlag{Name: "message", Aliases: []string{"m"}},
}

func main() {
	app := &cli.App{
		Name:   "changeset",
		Action: add.Run,
		Commands: []*cli.Command{
			{
				Name:   "add",
				Flags:  addFlags,
				Action: add.Run,
			},
			{
				Name:    "version",
				Aliases: []string{"consume"},
				Action:  version.Run,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "dry-run"},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
