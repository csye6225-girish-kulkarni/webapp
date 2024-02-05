package main

import (
	_ "github.com/lib/pq"
	"log"
	routerPkg "webapp/router"
)

func main() {
	router := routerPkg.InitializeRouter()
	err := router.Run(":8080")
	if err != nil {
		log.Printf("apperror starting server: %v", err)
		return
	}
}
