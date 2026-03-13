package main

import (
	"os"
	"shifumi/docs"
	"shifumi/internal/routes"
	"time"

	"github.com/gin-contrib/cors"
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
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}

	router := gin.Default()
	config := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
	router.Use(cors.New(config))
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
