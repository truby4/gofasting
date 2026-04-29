package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/truby4/gofasting/internal/app"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	env := flag.String("env", "dev", "environment")
	flag.Parse()

	var level slog.Level
	switch *env {
	case "dev":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	app, err := app.New(level)
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	err = app.Serve(*addr)
	app.Logger.Error(err.Error())
	os.Exit(1)
}
