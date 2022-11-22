package main

import (
	"log"

	"github.com/TsunamiProject/yamarkt/internal/config"
)

func main() {
	//Creating config instance
	log.Print("Initializing config")
	cfg := config.New()
	log.Print(cfg)
}
