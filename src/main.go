package main

import (
	"github.com/joho/godotenv"
	"log"
	"httpserver"
	"model"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	err := godotenv.Load("src/properties.env")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		genesisBlock := model.BuildGenesisBlock();
		spew.Dump(genesisBlock)
	}()
	log.Fatal(httpserver.Run())
}
