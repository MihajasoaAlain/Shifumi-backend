# Shifumi Backend

Backend API for a two-player Rock-Paper-Scissors game built with Go and Gin.

## Overview

This service provides:

- game creation
- player join flow
- round play handling
- live game updates over Server-Sent Events (SSE)
- Swagger documentation

The current implementation stores all games in memory. Data is lost when the process restarts.

## Stack

- Go `1.26`
- Gin
- Swagger via `swaggo/gin-swagger`
- SSE for live game events

## Project Structure

```text
.
├── cmd/server              # application entrypoint
├── internal/handlers       # HTTP handlers
├── internal/routes         # route registration
├── internal/services       # game logic and SSE broker
├── internal/models         # request/response and domain models
├── docs                    # generated Swagger files
├── Makefile
└── render.yaml             # Render deployment config
```

## Run Locally

### Requirements

- Go `1.26`
- optional: `swag` CLI if you want to regenerate Swagger docs

### Start the server

```bash
go run ./cmd/server
```

The API starts on `http://localhost:8080` by default.

You can also override the port:

```bash
PORT=9000 go run ./cmd/server
```

### Using Make

```bash
make run
```

`make run` regenerates Swagger docs first, so it expects the `swag` command to be installed.

## Swagger

Swagger UI is exposed at:

```text
http://localhost:8080/swagger/index.html
```

To regenerate the docs manually:

```bash
make swag
```

## API Endpoints

### Health check

```http
GET /health
```

Example response:

```json
{
  "message": "work well"
}
```

### Create a game

```http
POST /game
```

Example response:

```json
{
  "id": "game-1",
  "players": [],
  "status": "waiting"
}
```

### Join a game

```http
POST /game/:id/join
Content-Type: application/json
```

Request body:

```json
{
  "username": "alice"
}
```

When the second player joins, the game status becomes `ready`.

### Get game state

```http
GET /game/:id
```

### Play a round

```http
POST /game/:id/play
Content-Type: application/json
```

Request body:

```json
{
  "username": "alice",
  "choice": "rock"
}
```

Valid choices:

- `rock`
- `paper`
- `scissors`

Behavior:

- if only one player has submitted a choice, the API stores it and waits for the second player
- once both players have submitted, the API computes the winner, updates scores, resets round choices, and returns the result

### Stream live game events

```http
GET /game/:id/events
```

This endpoint opens an SSE stream and sends:

- an initial `game.snapshot` event with the current game state
- `game.updated` when players join or submit a choice
- `round.completed` when a round finishes

Example using `curl`:

```bash
curl -N http://localhost:8080/game/game-1/events
```

## Quick Test Flow

### 1. Create a game

```bash
curl -X POST http://localhost:8080/game
```

### 2. Join with player 1

```bash
curl -X POST http://localhost:8080/game/game-1/join \
  -H "Content-Type: application/json" \
  -d '{"username":"alice"}'
```

### 3. Join with player 2

```bash
curl -X POST http://localhost:8080/game/game-1/join \
  -H "Content-Type: application/json" \
  -d '{"username":"bob"}'
```

### 4. Play choices

```bash
curl -X POST http://localhost:8080/game/game-1/play \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","choice":"rock"}'
```

```bash
curl -X POST http://localhost:8080/game/game-1/play \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","choice":"scissors"}'
```

## Game Rules and Constraints

- a game supports at most 2 players
- usernames must be non-empty
- usernames must be unique within a game
- a player can submit only one choice per round
- rounds can be played only when 2 players have joined

## Deployment

The repository includes [`render.yaml`](/run/media/yuna/f4bd4336-3da5-48ad-87e4-db4d401d9671/PROJECT/Shifumi-backend/render.yaml) for Render deployment.

Configured behavior:

- build command: `go build -tags netgo -ldflags '-s -w' -o app ./cmd/server`
- start command: `./app`
- health check path: `/health`

## Notes

- CORS is currently open to all origins
- game state is process-local and not shared across instances
- there are no automated tests in the repository at the moment
