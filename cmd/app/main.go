package main

import (
	"flag"
	"log"
	"os"

	"github.com/truby4/go-fasting/internal/app"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	err = app.Serve(*addr)
	app.Logger.Error(err.Error())
	os.Exit(1)
}
