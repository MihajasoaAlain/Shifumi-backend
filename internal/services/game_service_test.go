package services

import (
	"testing"

	"shifumi/internal/models"
)

func resetState() {
	gamesMu.Lock()
	games = map[string]*models.Game{}
	gameCounter = 1
	eventBroker = newGameEventBroker()
	gamesMu.Unlock()
}

func TestDetermineWinner(t *testing.T) {
	cases := []struct {
		name           string
		first, second  models.Choice
		expectedWinner int
	}{
		{"draw rock", models.Rock, models.Rock, 0},
		{"draw paper", models.Paper, models.Paper, 0},
		{"rock beats scissors", models.Rock, models.Scissors, 1},
		{"paper beats rock", models.Paper, models.Rock, 1},
		{"scissors beats paper", models.Scissors, models.Paper, 1},
		{"scissors loses to rock", models.Scissors, models.Rock, 2},
		{"rock loses to paper", models.Rock, models.Paper, 2},
		{"paper loses to scissors", models.Paper, models.Scissors, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := determineWinner(tc.first, tc.second); got != tc.expectedWinner {
				t.Fatalf("determineWinner(%q, %q) = %d, want %d", tc.first, tc.second, got, tc.expectedWinner)
			}
		})
	}
}

func TestIsValidChoice(t *testing.T) {
	for _, c := range []models.Choice{models.Rock, models.Paper, models.Scissors} {
		if !isValidChoice(c) {
			t.Errorf("expected %q to be valid", c)
		}
	}
	for _, c := range []models.Choice{"", "lizard", "ROCK"} {
		if isValidChoice(c) {
			t.Errorf("expected %q to be invalid", c)
		}
	}
}

func TestCreateAndJoin(t *testing.T) {
	resetState()

	game := CreateGame()
	if game.ID != "game-1" {
		t.Fatalf("expected first game ID game-1, got %q", game.ID)
	}
	if game.Status != models.Waiting {
		t.Fatalf("expected status waiting, got %q", game.Status)
	}

	if _, err := JoinGame(game.ID, ""); err == nil {
		t.Error("expected error when joining with empty username")
	}
	if _, err := JoinGame("game-999", "alice"); err == nil {
		t.Error("expected error when joining unknown game")
	}

	updated, err := JoinGame(game.ID, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != models.Waiting {
		t.Fatalf("expected waiting with one player, got %q", updated.Status)
	}

	if _, err := JoinGame(game.ID, "alice"); err == nil {
		t.Error("expected error when username already taken")
	}

	updated, err = JoinGame(game.ID, "bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != models.Ready {
		t.Fatalf("expected ready with two players, got %q", updated.Status)
	}

	if _, err := JoinGame(game.ID, "carol"); err == nil {
		t.Error("expected error when game is full")
	}
}

func TestPlayRoundFullFlow(t *testing.T) {
	resetState()

	game := CreateGame()
	mustJoin(t, game.ID, "alice")
	mustJoin(t, game.ID, "bob")

	res, err := PlayRound(game.ID, "alice", models.Rock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["message"] != "choice saved, waiting for the other player" {
		t.Fatalf("expected waiting message, got %v", res["message"])
	}

	// The choice must never leak through the cloned game sent over the wire.
	clone, _ := GetGameByGame(game.ID)
	for _, p := range clone.Players {
		if p.Choice != "" {
			t.Errorf("clone leaked choice for %s: %q", p.Username, p.Choice)
		}
	}
	if !clone.Players[0].HasChosen || clone.Players[1].HasChosen {
		t.Errorf("expected only alice to have chosen, got %+v", clone.Players)
	}

	if _, err := PlayRound(game.ID, "alice", models.Paper); err == nil {
		t.Error("expected error when playing twice in the same round")
	}

	res, err = PlayRound(game.ID, "bob", models.Scissors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["result"] != "win" || res["winner"] != "alice" {
		t.Fatalf("expected alice to win, got %v", res)
	}

	scores := res["scores"].(map[string]int)
	if scores["alice"] != 1 || scores["bob"] != 0 {
		t.Fatalf("unexpected scores: %v", scores)
	}

	// After a completed round, choices reset and status returns to ready.
	final, _ := GetGameByGame(game.ID)
	if final.Status != models.Ready {
		t.Fatalf("expected ready after round, got %q", final.Status)
	}
	for _, p := range final.Players {
		if p.HasChosen {
			t.Errorf("expected choices reset after round, %s still flagged", p.Username)
		}
	}
}

func TestPlayRoundDraw(t *testing.T) {
	resetState()

	game := CreateGame()
	mustJoin(t, game.ID, "alice")
	mustJoin(t, game.ID, "bob")

	if _, err := PlayRound(game.ID, "alice", models.Rock); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res, err := PlayRound(game.ID, "bob", models.Rock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["result"] != "draw" {
		t.Fatalf("expected draw, got %v", res)
	}
	scores := res["scores"].(map[string]int)
	if scores["alice"] != 0 || scores["bob"] != 0 {
		t.Fatalf("draw should not change scores: %v", scores)
	}
}

func TestPlayRoundErrors(t *testing.T) {
	resetState()

	if _, err := PlayRound("game-999", "alice", models.Rock); err == nil {
		t.Error("expected error for unknown game")
	}

	game := CreateGame()
	mustJoin(t, game.ID, "alice")

	if _, err := PlayRound(game.ID, "alice", models.Rock); err == nil {
		t.Error("expected error when game is not full")
	}

	mustJoin(t, game.ID, "bob")

	if _, err := PlayRound(game.ID, "alice", "lizard"); err == nil {
		t.Error("expected error for invalid choice")
	}
	if _, err := PlayRound(game.ID, "carol", models.Rock); err == nil {
		t.Error("expected error for player not in game")
	}
}

func mustJoin(t *testing.T, gameID, username string) {
	t.Helper()
	if _, err := JoinGame(gameID, username); err != nil {
		t.Fatalf("join %s failed: %v", username, err)
	}
}
