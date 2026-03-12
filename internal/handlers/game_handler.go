package handlers

import (
	"net/http"
	"shifumi/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateGameHandler godoc
// @Summary      Créer une nouvelle partie
// @Description  Crée une nouvelle partie de Shifumi et renvoie les détails de la partie créée.
// @Tags         Game
// @Produce      json
// @Success      200
// @Router       /game [post]
func CreateGameHandler(c *gin.Context) {
	game := services.CreateGame()
	c.JSON(http.StatusOK, game)

}
