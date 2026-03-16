package models

type GameEvent struct {
	Type string      `json:"type"`
	Game *Game       `json:"game,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
