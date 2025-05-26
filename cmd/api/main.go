package main

import (
	"log"

	"github.com/rubenbupe/recipe-video-parser/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
