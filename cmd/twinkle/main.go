package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Crtrpt/twinkle"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	twinkle.InitFlag()
	ctx := context.Background()
	app := twinkle.NewApp(ctx)
	app.Run(ctx)
	_ = <-c
	app.Stop(ctx)
}
