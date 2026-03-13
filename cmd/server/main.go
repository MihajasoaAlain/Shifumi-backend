package main

import (
	"os"
	"shifumi/internal/routes"

	"github.com/gin-gonic/gin"
)

//@title Shifumi API
//@version 1.0
//@description API pour le jeu Shifumi (Pierre-Papier-Ciseaux).
//@contact.name Mihajasoa
//@contact.email mihajasoaalain85@gmail.com

//@host localhost:8080
//@BasePath /

func main() {
	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
