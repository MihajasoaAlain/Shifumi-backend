package handlers

import (
	"fmt"
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

// PlayGameHandler godoc
// @Summary      Jouer un tour de Shifumi
// @Description  Permet à un joueur de jouer un tour de Shifumi en fournissant son choix (pierre, papier ou ciseaux).
// @Tags         Game
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID de la partie"
// @Param        body body      models.PlayRequest true "Requête pour jouer un tour"
// @Success      200  {object}  models.Game
// @Failure      400  {object}  map[string]string
// @Router       /game/{id}/play [post]
func PlayGameHandler(c *gin.Context) {
	gameID := c.Param("id")
	var req models.PlayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	game, err := services.PlayRound(gameID, req.Username, req.Choice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}

// StreamGameEventsHandler godoc
// @Summary      Suivre une partie en temps reel
// @Description  Ouvre un flux SSE pour recevoir les mises a jour d'une partie.
// @Tags         Game
// @Produce      text/event-stream
// @Param        id   path      string  true  "ID de la partie"
// @Success      200  {string}  string  "SSE stream"
// @Failure      400  {object}  map[string]string
// @Router       /game/{id}/events [get]
func StreamGameEventsHandler(c *gin.Context) {
	gameID := c.Param("id")
	events, snapshot, err := services.SubscribeToGameEvents(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer services.UnsubscribeFromGameEvents(gameID, events)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	if err := writeSSEvent(c, "game.snapshot", models.GameEvent{
		Type: "game.snapshot",
		Game: snapshot,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			if err := writeSSEvent(c, event.Type, event); err != nil {
				return
			}
		}
	}
}

func writeSSEvent(c *gin.Context, name string, payload interface{}) error {
	c.SSEvent(name, payload)

	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
		return nil
	}

	return fmt.Errorf("streaming unsupported")
}
