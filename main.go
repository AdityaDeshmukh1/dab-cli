package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/adityadeshmukh1/dab-cli/cmd"
)

func main() {
	app := &cli.App {
		Name: "dab-cli",
		Usage: "A CLI music player powered by DAB API",
		Commands: []*cli.Command {
			cmd.LoginCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
