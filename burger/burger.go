package burger

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand/v2"
)

const (
	imageBasePath = "resources/images/ingredients/png/"
)

// Ingredient defines common properties of a particular Burger Part, e.g. the name and image of a tomato or cheese layer.
type Ingredient struct {
	name  string
	image *ebiten.Image
}

var allIngredients []Ingredient

var plateImage *ebiten.Image

func init() {
	allIngredients = []Ingredient{
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

	for index, ingredient := range allIngredients {
		image, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("%s/%s.png", imageBasePath, ingredient.name))

		if err != nil {
			log.Fatalf("Error while loading ingredient image %s: %v", ingredient.name, err)
		}

		allIngredients[index].image = image
	}

	var err error
	plateImage, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("%s/plate.png", imageBasePath))

	if err != nil {
		log.Fatalf("Error while loading plate image: %v", err)
	}
}

// Plate wraps an image object where burgers can be stacked on.
type Plate struct {
	image       *ebiten.Image
	x           int
	y           int
	scaleFactor float64
}

func (p *Plate) height() int {
	// Account for the scale factor when calculating the image height.
	return int(float64(p.image.Bounds().Size().Y) * p.scaleFactor)
}

func (p *Plate) draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(p.scaleFactor, p.scaleFactor)
	opts.GeoM.Translate(float64(p.x), float64(p.y))

	screen.DrawImage(p.image, opts)
}

// Part is a specific instances of an ingredient with a specific position on the screen. Note that there can be multiple
// Parts of the same Ingredient, but at different locations on the screen, e.g. multiple pickle instances.
type Part struct {
	ingredient  Ingredient
	lane        int
	x           int
	y           int
	scaleFactor float64
}

func newPart(ingredient Ingredient, lane int, scaleFactor float64) *Part {
	return &Part{
		ingredient:  ingredient,
		lane:        lane,
		x:           ScreenBorderSize + lane*laneWidth,
		y:           ScreenBorderSize,
		scaleFactor: scaleFactor,
	}
}

func newRandomFallingPart() *Part {
	if len(allIngredients) == 0 {
		log.Fatal("Ingredients not initialized yet, cannot create random part")
	}
	randomIngredient := allIngredients[rand.IntN(len(allIngredients))]

	// Spawn the Part in the middle lane.
	lane := laneCount / 2

	return newPart(randomIngredient, lane, buildScaleFactor)
}

func (p *Part) move(x, y int) {
	p.x = x
	p.y = y
}

func (p *Part) height() int {
	// Account for the scale factor when calculating the images height.
	return int(float64(p.ingredient.image.Bounds().Size().Y) * p.scaleFactor)
}

func (p *Part) draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(p.scaleFactor, p.scaleFactor)
	opts.GeoM.Translate(float64(p.x), float64(p.y))

	screen.DrawImage(p.ingredient.image, opts)
}

// Burger is a stack of Part instances on a Plate, i.e. specific ingredients at a certain location ordered from bottom to top.
type Burger struct {
	stack []*Part
	plate *Plate
}

func newEmptyBurger(lane int) *Burger {
	plate := &Plate{
		image:       plateImage,
		x:           ScreenBorderSize + lane*laneWidth,
		scaleFactor: buildScaleFactor,
	}
	plate.y = BuildSectionHeight - plate.height() - 1

	return &Burger{
		stack: make([]*Part, 0),
		plate: plate,
	}
}

func newRandomBurger(partCount int, lane int) *Burger {
	if len(allIngredients) == 0 {
		log.Fatal("Ingredients not initialized yet, cannot create random part")
	}

	stack := make([]*Part, partCount)

	// Always start with a bottom bun and finsh with a top bun.
	stack[0] = newPart(allIngredients[0], lane, orderScaleFactor)
	stack[partCount-1] = newPart(allIngredients[len(allIngredients)-1], lane, orderScaleFactor)

	for index := 1; index < partCount-1; index++ {
		// Pick a random ingredient, avoid repeating a bottom or top bun.
		randomIngredient := allIngredients[rand.IntN(len(allIngredients)-2)+1]

		stack[index] = newPart(randomIngredient, lane, orderScaleFactor)
	}

	plate := &Plate{
		image:       plateImage,
		x:           ScreenBorderSize + lane*laneWidth,
		scaleFactor: orderScaleFactor,
	}

	// Set y position starting from the plate, then bottom ingredient all the way up to the top.
	currentY := ScreenHeight - ScreenBorderSize - plate.height() - 1
	plate.y = currentY

	// Allow a large overlay for the first ingredient on the plate.
	currentY -= plate.height() / 5

	for _, part := range stack {
		part.y = currentY
		// Allow some overlay for the ingredients.
		currentY -= part.height() / 2
	}

	return &Burger{
		stack: stack,
		plate: plate,
	}
}

func (b *Burger) top() *Part {
	if len(b.stack) == 0 {
		return nil
	}

	return b.stack[len(b.stack)-1]
}

func (b *Burger) add(part *Part) {
	b.stack = append(b.stack, part)
}

func (b *Burger) draw(screen *ebiten.Image) {
	b.plate.draw(screen)

	for _, part := range b.stack {
		part.draw(screen)
	}
}
