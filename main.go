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
	width         = 1200
	height        = 1200
	stepPerTick   = 5
	scaleFactor   = 0.4
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
		{name: "tomatoes"},
		{name: "patty_beef"},
		{name: "ketchup"},
		{name: "salad"},
		{name: "patty_vegan"},
		{name: "mayo"},
		{name: "cheese"},
		{name: "ham"},
		{name: "onions"},
		{name: "pickles"},
		{name: "bun_top"},
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
	if g.falling == nil {
		// Piling finished, nothing to update.
		return nil
	}

	// By default, the ingredient should stop moving before at the lower bottom.
	compHeight := int(float64(g.falling.compType.image.Bounds().Size().Y) * scaleFactor)
	maxY := height - compHeight - 1

	if len(g.pile) > 0 {
		// Stop falling when we are halfway over the top most ingredient on the pile (allow some overlay).
		maxY = g.pile[len(g.pile)-1].y - (compHeight / 2)
	}
	//log.Printf("maxY: %d", maxY)

	if maxY < 0 || len(g.pile) >= len(allComponentTypes) {
		// Finish the piling process.
		//g.pile = make([]*component, 0)
		g.falling = nil

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

func (g *Game) Draw(screen *ebiten.Image) {
	drawComponent(screen, g.falling)

	for _, comp := range g.pile {
		drawComponent(screen, comp)
	}

	//ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, %d) (%f)", g.falling.x, g.falling.y, ebiten.ActualTPS()))
	//opts.GeoM.Scale(0.2, 0.2)
	//size := g.fruit.Bounds().Size()
	//x, y := ebiten.CursorPosition()
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("(%d, %d) - (%d, %d) - (%d, %d)", x, y, g.x, g.y, size.X, size.Y))
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", ebiten.ActualTPS()), 0, 10)
}

func drawComponent(screen *ebiten.Image, comp *component) {
	if comp == nil {
		// Nothing to do.
		return
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleFactor, scaleFactor)
	opts.GeoM.Translate(float64(comp.x), float64(comp.y))

	screen.DrawImage(comp.compType.image, opts)
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Burger Stash")

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatalf("Error while running game loop: %v", err)
	}
}
