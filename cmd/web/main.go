package main

import (
	"flag"
	"os"

	"github.com/truby4/go-fasting/internal/app"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	app := app.New()

	err := app.Serve(*addr)
	app.Logger.Error(err.Error())
	os.Exit(1)
}
