package services

import (
	"fmt"
	"shifumi/internal/models"
)

var games = map[string]*models.Game{}
var gameCounter = 1

func CreateGame() *models.Game {
	gameID := fmt.Sprintf("game-%d", gameCounter)
	gameCounter++
	game := &models.Game{
		ID:      gameID,
		Players: []models.Player{},
		Status:  "waiting",
	}
	games[gameID] = game
	return game

}
