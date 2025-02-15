package main

import (
	"log"

	"github.com/Prajna1999/atlas-be/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to inititalize application: %v", err)
		return
	}
	if err := application.Run(); err != nil {
		log.Fatalf("failed to run application %v", err)
		return
	}

}
