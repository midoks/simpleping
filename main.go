package main

import (
	// "fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"simpleping/cmd"
)

// Init config
var Version = "1.0"


func main() {

	app := cli.NewApp()
	app.Name = "simpleping"
	app.Usage = "simpleping service"
	app.Version = Version
	app.Commands = []cli.Command{
		cmd.Ping,
		cmd.Service,
	}

	if err := app.Run(os.Args); err != nil {
		log.Printf("Failed to start application: %v", err)
	}
}