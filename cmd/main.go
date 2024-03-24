package main

import (
	"log"

	"github.com/akshay0074700747/project-and-company-management-chat-service/config"
	"github.com/akshay0074700747/project-and-company-management-chat-service/initializer"
)

func main() {
	cfg, err := config.LoadConfigurations()
	if err != nil {
		log.Fatal(err, "cannot load envs")
	}

	handler := initializer.Inject(&cfg)
	handler.Start()     
}
