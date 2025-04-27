package main

import (
	"log"

	"github.com/KrllF/Cloud/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Println("app.NewApp: ", err)

		return
	}
	if err := a.Run(); err != nil {
		log.Panic("a.Run(): ", err.Error())
	}
}
