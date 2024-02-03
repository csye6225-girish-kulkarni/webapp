package main

import (
	routerPkg "Health-Check/router"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	router := routerPkg.InitializeRouter()
	err := router.Run(":8080")
	if err != nil {
		log.Printf("error starting server: %v", err)
		return
	}
}
