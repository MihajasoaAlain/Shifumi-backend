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

func JoinGame(gameID string, username string) (*models.Game, error) {
	game, exists := games[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if len(game.Players) >= 2 {
		return nil, fmt.Errorf("game is full")
	}
	for _, player := range game.Players {
		if player.Username == username {
			return nil, fmt.Errorf("username already taken in this game")
		}
	}
	newPlayer := models.Player{
		Username: username,
		Choice:   "",
		Score:    0,
	}
	game.Players = append(game.Players, newPlayer)
	if len(game.Players) == 2 {
		game.Status = "ready"
	}
	return game, nil
}

func GetGameByGame(gameID string) (*models.Game, error) {
	game, exists := games[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}
