package main

import (
	"context"
	"http_bridge"
)

func main() {
	do := make(chan struct{}, 0)
	http_bridge.InitFlag()
	ctx := context.Background()
	app := http_bridge.NewApp(ctx)
	app.Run(ctx)
	<-do
}
