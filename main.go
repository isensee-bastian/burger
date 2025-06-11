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

// componentType holds immutable properties of a particular burger component.
type componentType struct {
	name  string
	index int
	image *ebiten.Image
}

var allComponentTypes []componentType

type component struct {
	compType componentType
	x        int
	y        int
}

func newComponent(compType componentType) *component {
	return &component{compType: compType}
}

func init() {
	allComponentTypes = []componentType{
		{name: "bun_bottom"},
		{name: "bun_top"},
		{name: "cheese"},
		{name: "ham"},
		{name: "ketchup"},
		{name: "mayo"},
		{name: "onions"},
		{name: "patty"},
		{name: "salad"},
		{name: "tomatoes"},
	}

	for index, compType := range allComponentTypes {
		path := fmt.Sprintf("%s/%s.png", imageBasePath, compType.name)
		image, _, err := ebitenutil.NewImageFromFile(path)

		if err != nil {
			log.Fatalf("Error while loading image %s: %v", compType.name, err)
		}

		allComponentTypes[index].index = index
		allComponentTypes[index].image = image
	}
}

type Game struct {
	falling *component
	pile    []*component
}

func newGame() *Game {
	return &Game{
		falling: newComponent(allComponentTypes[0]),
		pile:    make([]*component, 0),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	compSize := g.falling.compType.image.Bounds().Size()
	maxY := height - compSize.Y - 1

	if len(g.pile) > 0 {
		maxY = g.pile[len(g.pile)-1].y - (compSize.Y / 2)
	}
	//log.Printf("maxY: %d", maxY)

	if maxY < 0 {
		// Restart the whole piling process.
		g.pile = make([]*component, 0)

		return nil
	}

	g.falling.y += stepPerTick

	if g.falling.y > maxY {
		g.falling.y = maxY
		g.pile = append(g.pile, g.falling)

		nextComponentIndex := (g.falling.compType.index + 1) % len(allComponentTypes)
		g.falling = newComponent(allComponentTypes[nextComponentIndex])
	}
	//log.Printf("falling.y: %d", g.falling.y)

	return nil
}

func drawComponent(screen *ebiten.Image, comp *component) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(comp.x), float64(comp.y))
	screen.DrawImage(comp.compType.image, opts)
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawComponent(screen, g.falling)

	for _, comp := range g.pile {
		drawComponent(screen, comp)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, %d) (%f)", g.falling.x, g.falling.y, ebiten.ActualTPS()))
	//opts.GeoM.Scale(0.2, 0.2)
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
