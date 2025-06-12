package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"log"
	"math/rand/v2"
)

const (
	width         = 1200
	height        = 1200
	stepPerTick   = 5
	scaleFactor   = 0.4
	laneCount     = 3
	laneWidth     = width / laneCount
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
	lane     int
	x        int
	y        int
}

func newComponent(compType componentType) *component {
	// Spawn the component in the middle lane.
	lane := laneCount / 2

	return &component{
		compType: compType,
		lane:     lane,
		x:        lane * laneWidth,
		y:        0,
	}
}

func newRandomComponent() *component {
	if len(allComponentTypes) == 0 {
		log.Fatal("Component types not initialized yet, cannot create random component")
	}

	randomType := allComponentTypes[rand.IntN(len(allComponentTypes))]
	return newComponent(randomType)
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
	// falling is the currently moving ingredient that needs to be steered onto a pile.
	falling *component
	// lanes represents multiple piles of layered ingredient, from left to right.
	lanes []*pile
}

type pile struct {
	// components contains layers of ingredients from bottom (index 0) to top (index len - 1).
	components []*component
}

func newGame() *Game {
	lanes := make([]*pile, laneCount)
	for index := range laneCount {
		lanes[index] = &pile{components: make([]*component, 0)}
	}

	return &Game{
		falling: newComponent(allComponentTypes[0]),
		lanes:   lanes,
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
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.falling.lane = max(0, g.falling.lane-1)
		g.falling.x = g.falling.lane * laneWidth
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.falling.lane = min(laneCount-1, g.falling.lane+1)
		g.falling.x = g.falling.lane * laneWidth
	}

	// By default, the ingredient should stop moving before at the lower bottom.
	compHeight := int(float64(g.falling.compType.image.Bounds().Size().Y) * scaleFactor)
	maxY := height - compHeight - 1

	currentPile := g.lanes[g.falling.lane]

	if len(currentPile.components) > 0 {
		// Stop falling when we are halfway over the top most ingredient on the pile (allow some overlay).
		maxY = currentPile.components[len(currentPile.components)-1].y - (compHeight / 2)
	}
	//log.Printf("maxY: %d", maxY)

	if maxY < 0 {
		// Finish the piling process.
		g.falling = nil

		return nil
	}

	g.falling.y += stepPerTick

	if g.falling.y > maxY {
		g.falling.y = maxY
		currentPile.components = append(currentPile.components, g.falling)

		g.falling = newRandomComponent()
	}
	//log.Printf("falling.y: %d", g.falling.y)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawComponent(screen, g.falling)

	for _, lane := range g.lanes {
		for _, comp := range lane.components {
			drawComponent(screen, comp)
		}
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
