package burger

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

const (
	ScreenWidth  = 1200
	ScreenHeight = 1200

	defaultFallStep = 5
	fastFallStep    = defaultFallStep * 4
	scaleFactor     = 0.4
	laneCount       = 3
	laneWidth       = ScreenWidth / laneCount

	imageBasePath = "resources/images/ingredients/png/"
)

type Game struct {
	// falling is the currently moving ingredient Part that has not yet landed on a specific Burger.
	falling *Part
	// burgers represents multiple stacks of layered ingredient Parts, ordered from left to right (lanes).
	burgers []*Burger
}

func (g *Game) sellAllowed(burgerIndex int) bool {
	if burgerIndex < 0 || burgerIndex >= len(g.burgers) {
		return false
	}

	return g.burgers[burgerIndex].top() != nil
}

func (g *Game) sell(burgerIndex int) {
	g.burgers[burgerIndex] = newBurger()
}

func (g *Game) move(targetLane int, stepSize int) {
	if g.falling == nil {
		// Piling finished, nothing to update.
		return
	}

	if targetLane < 0 || targetLane >= len(g.burgers) {
		// Refuse changing the lane as we are at the bounds of the screen.
		targetLane = g.falling.lane
	}

	// By default, if there is no started burger yet, the ingredient part should stop moving at the lower bottom.
	partHeight := g.falling.height()
	maxY := ScreenHeight - partHeight - 1

	targetBurger := g.burgers[targetLane]
	topPart := targetBurger.top()

	if topPart != nil {
		// There is an existing burger, stop moving when we are halfway over the top most part (allow some overlay).
		maxY = topPart.y - (partHeight / 2)
	}

	if maxY < 0 {
		// Finish the piling and spawning process.
		g.falling = nil

		return
	}

	if g.falling.y > maxY {
		// There is no more space in this Burgers lane to add new Parts. We need to prevent a collision.
		// Refuse changing the lane and move normally if a vertical move was initiated.
		// Else finish the piling and spawning process.
		if targetLane == g.falling.lane {
			// Finish the piling and spawning process.
			g.falling = nil

			return
		}

		// TODO: Check if we can avoid recursive calls here and make it more simple.
		g.move(g.falling.lane, stepSize)

		return
	}

	// Accept the lane change if present as we found no collision with existing Burgers.
	g.falling.lane = targetLane

	// Move our falling Part one step further down as we have not hit the bottom or Burger top yet.
	g.falling.move(g.falling.lane*laneWidth, g.falling.y+stepSize)

	if g.falling.y > maxY {
		// We hit the bottom or the top most Part of a Burger. Add the layer to the burger and spawn a new Part.
		g.falling.move(g.falling.x, maxY)
		targetBurger.add(g.falling)
		g.falling = newRandomPart()
	}
}

func NewGame() *Game {
	burgers := make([]*Burger, laneCount)

	for index := range laneCount {
		burgers[index] = newBurger()
	}

	return &Game{
		falling: newRandomPart(),
		burgers: burgers,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally.
		return ebiten.Termination
	}
	if g.falling == nil {
		// Piling finished, nothing to update.
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.move(g.falling.lane-1, defaultFallStep)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.move(g.falling.lane+1, defaultFallStep)
	} else if inpututil.KeyPressDuration(ebiten.KeyDown) > 0 {
		g.move(g.falling.lane, fastFallStep)
	} else if inpututil.IsKeyJustPressed(ebiten.Key1) && g.sellAllowed(0) {
		g.sell(0)
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) && g.sellAllowed(1) {
		g.sell(1)
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) && g.sellAllowed(2) {
		g.sell(2)
	} else {
		g.move(g.falling.lane, defaultFallStep)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.falling != nil {
		g.falling.draw(screen)
	}

	for _, burger := range g.burgers {
		burger.draw(screen)
	}
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}
