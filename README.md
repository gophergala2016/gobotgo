# GoBotGo

GoBotGo is a competitive Go (board game) arena written in Go (golang) for bots to play via an HTTP API. Hosted at http://gobotgo.bellstone.ca.

## Premise

- The pun is good.
- We wanted to be able to compare bots written in different languages.

## What Works

- API endpoint.
- Board state, game rules and point totals.
- Go client library for writing bots.
- Demo bots, both random and best available move.
- Sketchy Human-AI/Human-Human interface.

## API

- All requests are under the root `/api/v1/game/`.
- `start/` returns a GameID and starting color. Each player is given their own GameID.
- `play/<GameID>/move` accepts a move as `[]` (pass) or `[x, y]`.
- `play/<GameID>/state` returns board state, including current player, captured and remaining pieces.
- `play/<GameID>/wait` returns after opponent has finished their turn.

## Todo

- Implement a game room.
- Register players and bots to compare histories.
- Statistics collection, game persistence.
- Allow trials to be setup to compare bots.
- Clean up javascript errors (lol).
- Write more bots!

## Known Bugs
- `plate/<GameID>/state` returns game over once the game is over.