package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var (
	checkURL = ""
)

func main() {
	app := cli.NewApp()

	app.Name = "Http-go-checker"
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		{
			Name:  "Arturo Reuschenbach Puncernau",
			Email: "a.reuschenbach.puncernau@sap.com",
		},
	}
	app.Usage = "check http connections"
	app.Action = runChecker
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url,U",
			Usage: "url to check",
			Value: "www.google.com",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// private

func runChecker(c *cli.Context) {
	if c.GlobalString("url") != "" {
		checkURL = c.GlobalString("url")
	} else {
		log.Fatalf("Url not provided")
	}

	log.Infof("Checking URL: %s", checkURL)

}
