package main

import (
	"log"
	"os"

	"github.com/chyroc/write-deprecated/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "write-deprecated",
		UsageText: "Write deprecated comment to go source file",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "dir", Usage: "write which dir"},
			&cli.StringFlag{Name: "comment", Usage: "comment for deprecated func/type"},
		},
		Action: func(c *cli.Context) error {
			dir := c.String("dir")
			comment := c.String("comment")

			return internal.WriteDeprecated(dir, comment)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
