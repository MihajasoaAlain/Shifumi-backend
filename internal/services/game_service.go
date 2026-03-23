package services

import (
	"fmt"
	"sync"

	"shifumi/internal/models"
)

var (
	games       = map[string]*models.Game{}
	gameCounter = 1
	gamesMu     sync.RWMutex
	eventBroker = newGameEventBroker()
)

func CreateGame() *models.Game {
	gamesMu.Lock()
	defer gamesMu.Unlock()

	gameID := fmt.Sprintf("game-%d", gameCounter)
	gameCounter++
	game := &models.Game{
		ID:      gameID,
		Players: []models.Player{},
		Status:  models.Waiting,
	}
	games[gameID] = game

	eventBroker.Publish(gameID, models.GameEvent{
		Type: "game.created",
		Game: cloneGame(game),
	})

	return game

}

func JoinGame(gameID string, username string) (*models.Game, error) {
	gamesMu.Lock()
	defer gamesMu.Unlock()

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

	eventBroker.Publish(gameID, models.GameEvent{
		Type: "game.updated",
		Game: cloneGame(game),
		Data: map[string]string{
			"action":   "player_joined",
			"username": username,
		},
	})

	return cloneGame(game), nil
}

func GetGameByGame(gameID string) (*models.Game, error) {
	gamesMu.RLock()
	defer gamesMu.RUnlock()

	game, exists := games[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return cloneGame(game), nil
}

func GetWaitingGames() ([]*models.Game, error) {
	gamesMu.RLock()
	defer gamesMu.RUnlock()

	waitingGames := make([]*models.Game, 0)

	for _, game := range games {
		if game.Status == models.Waiting {
			waitingGames = append(waitingGames, cloneGame(game))
		}
	}
	return waitingGames, nil
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
	gamesMu.Lock()
	defer gamesMu.Unlock()

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
		return nil, fmt.Errorf("player has already played this round")
	}

	game.Players[playerIndex].Choice = choice
	game.Status = models.Playing

	if game.Players[0].Choice == "" || game.Players[1].Choice == "" {
		result := map[string]interface{}{
			"message": "choice saved, waiting for the other player",
			"game":    cloneGame(game),
		}

		eventBroker.Publish(gameID, models.GameEvent{
			Type: "game.updated",
			Game: cloneGame(game),
			Data: map[string]string{
				"action":   "choice_submitted",
				"username": username,
			},
		})

		return result, nil
	}

	firstChoice := game.Players[0].Choice
	secondChoice := game.Players[1].Choice

	winner := determineWinner(firstChoice, secondChoice)

	var result map[string]interface{}

	switch winner {
	case 0:
		result = map[string]interface{}{
			"message": "round completed",
			"result":  "draw",
			"choices": map[string]models.Choice{
				game.Players[0].Username: firstChoice,
				game.Players[1].Username: secondChoice,
			},
			"scores": map[string]int{
				game.Players[0].Username: game.Players[0].Score,
				game.Players[1].Username: game.Players[1].Score,
			},
		}
	case 1:
		game.Players[0].Score++
		result = map[string]interface{}{
			"message": "round completed",
			"result":  "win",
			"winner":  game.Players[0].Username,
			"choices": map[string]models.Choice{
				game.Players[0].Username: firstChoice,
				game.Players[1].Username: secondChoice,
			},
			"scores": map[string]int{
				game.Players[0].Username: game.Players[0].Score,
				game.Players[1].Username: game.Players[1].Score,
			},
		}
	case 2:
		game.Players[1].Score++
		result = map[string]interface{}{
			"message": "round completed",
			"result":  "win",
			"winner":  game.Players[1].Username,
			"choices": map[string]models.Choice{
				game.Players[0].Username: firstChoice,
				game.Players[1].Username: secondChoice,
			},
			"scores": map[string]int{
				game.Players[0].Username: game.Players[0].Score,
				game.Players[1].Username: game.Players[1].Score,
			},
		}
	default:
		return nil, fmt.Errorf("unexpected winner value")
	}

	for i := range game.Players {
		game.Players[i].Choice = ""
	}

	game.Status = models.Ready

	eventBroker.Publish(gameID, models.GameEvent{
		Type: "round.completed",
		Game: cloneGame(game),
		Data: result,
	})

	return result, nil
}

func SubscribeToGameEvents(gameID string) (chan models.GameEvent, *models.Game, error) {
	gamesMu.RLock()
	game, exists := games[gameID]
	if !exists {
		gamesMu.RUnlock()
		return nil, nil, fmt.Errorf("game not found")
	}
	snapshot := cloneGame(game)
	gamesMu.RUnlock()

	return eventBroker.Subscribe(gameID), snapshot, nil
}

func UnsubscribeFromGameEvents(gameID string, ch chan models.GameEvent) {
	eventBroker.Unsubscribe(gameID, ch)
}

func cloneGame(game *models.Game) *models.Game {
	if game == nil {
		return nil
	}

	clonedPlayers := make([]models.Player, len(game.Players))
	copy(clonedPlayers, game.Players)

	return &models.Game{
		ID:      game.ID,
		Players: clonedPlayers,
		Status:  game.Status,
	}
}
