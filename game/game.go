// Package game provides means of playing a game
package game

type player int

const (
	Empty = player(iota)
	Black
	White
)

type Move struct {
	Color player
	X, Y  int
}
