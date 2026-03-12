package models

type Player struct {
	Username string `json:"username"`
	Choice   string `json:"choice"`
	Score    int    `json:"score"`
}

type Game struct {
	ID      string   `json:"id"`
	Players []Player `json:"players"`
	Status  string   `json:"status"`
}

type JoinGameRequest struct {
	Username string `json:"username"`
}
