package main

import (
	"github.com/alexflint/go-arg"

	_ "github.com/arthur-trt/bechdelproxy/database"
	"github.com/arthur-trt/bechdelproxy/log"
	"github.com/arthur-trt/bechdelproxy/movies"
	"github.com/arthur-trt/bechdelproxy/api"
)

type args struct {
	Update bool `help:"Update the local DB with https://bechdeltest.com/ DB"`
}

func (args) Version() string {
	return "v0.1.0"
}

func (args) Description() string {
	return "This program is proxy, it will serve ratings from https://bechdeltest.com/ with a local database."
}

func (args) Epilogue() string {
	return "For more information visit github.com/arthur-trt/bechdelproxy"
}

func main() {
	var arguments args
	arg.MustParse(&arguments)

	if arguments.Update {
		log.Info("Starting movie DB update")
		if err := movies.Update(); err != nil {
			log.Fatal("An error occurred while updating the database")
		}
		log.Info("Ending movie DB update")
	} else {
		echo := api.New()
		api.Register(echo)
		log.Fatal(echo.Start(":1789"))
	}
}
