package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
)

const (
	width         = 1000
	height        = 1000
	stepPerTick   = 5
	imageBasePath = "resources/images/ingredients/png/"
)

type ingredientType int

const (
	bunBottom ingredientType = iota
	bunTop
	cheese
	ham
	ketchup
	mayo
	onions
	patty
	salad
	tomatoes
)

type ingredient struct {
	name           string
	ingredientType ingredientType
	image          *ebiten.Image
}

var ingredients []ingredient

func init() {
	ingredients = []ingredient{
		{name: "bun_bottom", ingredientType: bunBottom},
		{name: "bun_top", ingredientType: bunTop},
		{name: "cheese", ingredientType: cheese},
		{name: "ham", ingredientType: ham},
		{name: "ketchup", ingredientType: ketchup},
		{name: "mayo", ingredientType: mayo},
		{name: "onions", ingredientType: onions},
		{name: "patty", ingredientType: patty},
		{name: "salad", ingredientType: salad},
		{name: "tomatoes", ingredientType: tomatoes},
	}

	for index, ingredient := range ingredients {
		path := fmt.Sprintf("%s/%s.png", imageBasePath, ingredient.name)
		image, _, err := ebitenutil.NewImageFromFile(path)

		if err != nil {
			log.Fatalf("Error while loading image %s: %v", ingredient.name, err)
		}

		ingredients[index].image = image
	}
}

type Game struct {
	ingredientIndex int
	x               int
	y               int
}

func newGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	g.y += stepPerTick

	if g.y >= height {
		g.y = 0
		g.ingredientIndex = (g.ingredientIndex + 1) % len(ingredients)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(g.x), float64(g.y))
	//opts.GeoM.Scale(0.2, 0.2)

	screen.DrawImage(ingredients[g.ingredientIndex].image, opts)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, %d) (%f)", g.x, g.y, ebiten.ActualTPS()))
	//size := g.fruit.Bounds().Size()
	//x, y := ebiten.CursorPosition()
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, %d) - (%d, %d) - (%d, %d)", x, y, g.x, g.y, size.X, size.Y))
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", ebiten.ActualTPS()), 0, 10)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Burger Challenge")

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatalf("Error while running game loop: %v", err)
	}
}
