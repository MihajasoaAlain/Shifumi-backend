package services

import (
	"errors"
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
		Status:  models.Waiting,
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
		game.Status = models.Ready
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

func isValidChoice(choice models.Choice) bool {
	return choice == models.Rock ||
		choice == models.Paper ||
		choice == models.Scissors
}

func determineWinner(choice1 models.Choice, choice2 models.Choice) int {
	if choice1 == choice2 {
		return 0
	}
	if (choice1 == models.Rock && choice2 == models.Scissors) ||
		(choice1 == models.Paper && choice2 == models.Rock) ||
		(choice1 == models.Scissors && choice2 == models.Paper) {
		return 1
	}
	return 2
}

func PlayRound(gameID string, username string, choice models.Choice) (map[string]interface{}, error) {
	fmt.Printf("PlayRound called with gameID: %s, username: %s, choice: %s\n", gameID, username, choice)
	game, exists := games[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	if len(game.Players) < 2 {
		return nil, fmt.Errorf("game is not full yet")
	}
	if game.Status != models.Playing && game.Status != models.Ready {
		return nil, fmt.Errorf("game is not in playing status")
	}

	if !isValidChoice(choice) {
		return nil, fmt.Errorf("invalid choice")
	}
	playerIndex := -1
	for i, player := range game.Players {
		if player.Username == username {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		return nil, fmt.Errorf("player not found in this game")
	}

	if game.Players[playerIndex].Choice != "" {
		return nil, errors.New("player has already played this round")
	}
	if game.Players[playerIndex].Choice != "" {
		return nil, errors.New("player has already played this round")
	}
	game.Players[playerIndex].Choice = choice
	game.Status = models.Playing

	if game.Players[0].Choice != "" && game.Players[1].Choice != "" {
		return map[string]interface{}{
			"message": "choice saved, waiting for the other player",
			"game":    game,
		}, nil
	}
	firstChoice := game.Players[0].Choice
	secondChoice := game.Players[1].Choice
	winner := determineWinner(firstChoice, secondChoice)
	switch winner {
	case 1:
		game.Players[0].Score++
	case 2:
		game.Players[1].Score++
	}

	result := map[string]interface{}{
		"message": "round completed",
		"game":    game,
	}
	for i := range game.Players {
		game.Players[i].Choice = ""
	}
	return result, nil
}
