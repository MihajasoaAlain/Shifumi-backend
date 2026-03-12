package models

type Player struct {
	Username string `json:"username"`
	Choice   Choice `json:"choice"`
	Score    int    `json:"score"`
}

type Game struct {
	ID      string     `json:"id"`
	Players []Player   `json:"players"`
	Status  GameStatus `json:"status"`
}

type JoinGameRequest struct {
	Username string `json:"username"`
}

type PlayRequest struct {
	Username string `json:"username"`
	Choice   Choice `json:"choice"`
}
