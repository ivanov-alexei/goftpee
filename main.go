package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ivanov-alexei/goftpee/server"
)

func main() {
	config := server.Config{}
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Println("Read file error")
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("Unmarshalling error ")
		panic(err)
	}

	s := server.NewFtpServer(&config)
	err = s.Start()
	if err != nil {
		panic(err)
	}
}
