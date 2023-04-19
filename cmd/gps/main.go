package main

import (
	"context"
	"fmt"
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
	s := <-c
	fmt.Printf("signal %s", s)
}
