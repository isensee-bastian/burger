package burger

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

const (
	ScreenWidth        = 1200
	ScreenHeight       = 1200
	BuildSectionHeight = 1000

	defaultFallStep  = 5
	fastFallStep     = defaultFallStep * 4
	buildScaleFactor = 0.4
	orderScaleFactor = 0.2
	laneCount        = 3
	laneWidth        = ScreenWidth / laneCount
)

type Game struct {
	// falling is the currently moving ingredient Part that has not yet landed on a specific Burger.
	falling *Part
	// burgers represents multiple stacks of layered ingredient Parts, ordered from left to right (lanes).
	burgers []*Burger
	// orders represents specific burger compositions requested by customers, ordered from left to right (lanes).
	orders []*Burger

	audioSell    *AudioPlayer
	audioStacked *AudioPlayer
}

func NewGame() *Game {
	burgers := make([]*Burger, laneCount)
	orders := make([]*Burger, laneCount)

	for lane := range laneCount {
		burgers[lane] = newEmptyBurger(lane)
		orders[lane] = newRandomBurger(7, lane)
	}

	return &Game{
		falling:      newRandomFallingPart(),
		burgers:      burgers,
		orders:       orders,
		audioSell:    newMp3AudioPlayer("cash_register.mp3"),
		audioStacked: newMp3AudioPlayer("plop.mp3"),
	}
}

func (g *Game) Close() {
	g.audioSell.Close()
	g.audioStacked.Close()
}

func (g *Game) sellAllowed(lane int) bool {
	if lane < 0 || lane >= len(g.burgers) {
		return false
	}

	return g.burgers[lane].top() != nil
}

func (g *Game) sell(lane int) {
	g.burgers[lane] = newEmptyBurger(lane)
	g.orders[lane] = newRandomBurger(7, lane)
	g.audioSell.Replay()
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

	targetBurger := g.burgers[targetLane]

	// By default, if there is no started burger yet, the ingredient part should stop moving when we hit the plate (allow large overlay).
	partHeight := g.falling.height()
	maxY := targetBurger.plate.y - (partHeight / 5)

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
		g.audioStacked.Replay()

		g.falling = newRandomFallingPart()
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) && g.sellAllowed(0) {
		g.sell(0)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) && g.sellAllowed(1) {
		g.sell(1)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) && g.sellAllowed(2) {
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

	for _, order := range g.orders {
		order.draw(screen)
	}
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}
