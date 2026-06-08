package main

import (
	"log"

	"gin-M-TIX/config"
	"gin-M-TIX/routes"
)

func main() {
	db := config.NewDatabase()
	router := routes.SetupRouter(db)

	if err := router.Run(":8999"); err != nil {
		log.Fatal(err)
	}
}
