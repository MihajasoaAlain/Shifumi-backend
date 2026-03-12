package models

type GameStatus string

const (
	Waiting GameStatus = "waiting"
	Ready   GameStatus = "ready"
	Playing GameStatus = "playing"
)
