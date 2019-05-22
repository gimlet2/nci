package main

import (
	"log"
	"nci/core"
)

func main() {
	log.Println("Application started")
	core.SetupRouting()
}
