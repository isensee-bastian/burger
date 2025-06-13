package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/isensee-bastian/burger/burger"
	"log"
)

func main() {
	ebiten.SetWindowSize(burger.ScreenWidth, burger.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Burger Stacker")

	game := burger.NewGame()
	defer game.Close()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Error while running game loop: %v", err)
	}
}
