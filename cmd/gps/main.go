package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Crtrpt/gps"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill)
	gps.InitFlag()
	ctx := context.Background()
	app := gps.NewApp(ctx)
	app.Run(ctx)
	_ = <-c
}
