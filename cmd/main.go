package main

import (
	"log"
	"os"
)

func main() {
	// create a log file
	f, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	// set log output to file and stdout
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Starting application...")
	log.Println("Initializing server...")
	log.Println("Server started successfully")

}
