package cmd

import (
	"log"
	"os"

	"github.com/headwindfly/go-auth0-web-app/internal/core"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		EnableBashCompletion: true,
		Version:              core.Version,
		Name:                 "go-auth0-web-app",
		Usage:                "Go Web Application",
		Before: func(c *cli.Context) error {
			if err := godotenv.Load(); err != nil {
				return err
			}

			return nil
		},
	}
)

// Execute executes commands.
func Execute() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
