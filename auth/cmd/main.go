package main

import (
	"flag"
	"log"

	"github.com/ilovepitsa/happy/auth/internal/app"
)

var (
	configPath = flag.String("c", "config/config.yaml", "use for custom config")
)

func main() {
	if err := app.Run(*configPath); err != nil {
		log.Fatal(err)
	}
}
