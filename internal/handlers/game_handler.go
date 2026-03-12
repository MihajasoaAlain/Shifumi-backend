package handlers

import (
	"net/http"
	"shifumi/internal/models"
	"shifumi/internal/services"

	_ "shifumi/docs"
	_ "shifumi/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateGameHandler godoc
// @Summary      Créer une nouvelle partie
// @Description  Crée une nouvelle partie de Shifumi et renvoie les détails de la partie créée.
// @Tags         Game
// @Produce      json
// @Success      200  {object}  models.Game
// @Router       /game [post]
func CreateGameHandler(c *gin.Context) {
	game := services.CreateGame()
	c.JSON(http.StatusOK, game)
}

// JoinGameHandler godoc
// @Summary      Rejoindre une partie existante
// @Description  Permet à un joueur de rejoindre une partie existante en fournissant un nom d'utilisateur.
// @Tags         Game
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID de la partie"
// @Param        body body      models.JoinGameRequest true "Requête pour rejoindre une partie"
// @Success      200  {object}  models.Game
// @Failure      400  {object}  map[string]string
// @Router       /game/{id}/join [post]
func JoinGameHandler(c *gin.Context) {
	gameID := c.Param("id")
	var req models.JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	game, err := services.JoinGame(gameID, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}

// GetGameHandler godoc
// @Summary      Obtenir les détails d'une partie
// @Description  Récupère les détails d'une partie de Shifumi en fonction de son ID.
// @Tags         Game
// @Produce      json
// @Param        id   path      string  true  "ID de la partie"
// @Success      200  {object}  models.Game
// @Failure      400  {object}  map[string]string
// @Router       /game/{id} [get]
func GetGameHandler(c *gin.Context) {
	gameID := c.Param("id")
	game, err := services.GetGameByGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}
