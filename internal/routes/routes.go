package routes

import (
	"shifumi/internal/handlers"

	_ "shifumi/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", handlers.HealthHandler)
	router.POST("/game", handlers.CreateGameHandler)
	router.POST("/game/:id/join", handlers.JoinGameHandler)
	router.GET("/game/:id", handlers.GetGameHandler)

}
