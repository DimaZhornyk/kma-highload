package main

import (
	"fmt"
	"log"
	"os"

	"lab4/service"

	"github.com/joho/godotenv"
)

func main()  {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	s, err := service.New()
	if err != nil {
		log.Fatalf("error bootstraping application: %s", err.Error())
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s.Echo().Logger.Fatal(s.Echo().Start(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))))
}
