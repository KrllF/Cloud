package main

import (
	"context"
	"log"

	"github.com/KrllF/Cloud/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Println("app.NewApp: ", err)

		return
	}
	if err := a.Run(); err != nil {
		log.Panic("a.Run(): ", err.Error())
	}
}
