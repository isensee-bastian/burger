package burger

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/isensee-bastian/burger/resources/images/ingredients/png"
	"image"
	"log"
	"math/rand/v2"
)

type IngredientType int

// Note that the IngredientType enum values MUST be kept ascending from 0 to N without any gaps for the
// randomization to work correctly. Moreover, BunBottom and BunTop must be kept at the beginning of the enum.
const (
	IngBunBottom IngredientType = iota
	IngBunTop
	IngPattyBeef
	IngPattyVegan
	IngHam
	IngTomatoes
	IngSalad
	IngKetchup
	IngMayo
	IngCheese
	IngOnions
	IngPickles
)

// Ingredient defines common properties of a particular Burger Part, e.g. the name and image of a tomato or cheese layer.
type Ingredient struct {
	typ   IngredientType
	image *ebiten.Image
}

var typeToIngredient = map[IngredientType]Ingredient{}

var plateImage *ebiten.Image

func init() {
	typeToRawImage := map[IngredientType][]byte{
		IngBunBottom:  png.BunBottom,
		IngBunTop:     png.BunTop,
		IngPattyBeef:  png.PattyBeef,
		IngPattyVegan: png.PattyVegan,
		IngHam:        png.Ham,
		IngTomatoes:   png.Tomatoes,
		IngSalad:      png.Salad,
		IngKetchup:    png.Ketchup,
		IngMayo:       png.Mayo,
		IngCheese:     png.Cheese,
		IngOnions:     png.Onions,
		IngPickles:    png.Pickles,
	}

	for typ, rawImage := range typeToRawImage {
		ingredientImage, _, err := image.Decode(bytes.NewReader(rawImage))

		if err != nil {
			log.Fatalf("Error while loading image for ingredient type %v: %v", typ, err)
		}

		typeToIngredient[typ] = Ingredient{
			typ:   typ,
			image: ebiten.NewImageFromImage(ingredientImage),
		}
	}

	plateStdImage, _, err := image.Decode(bytes.NewReader(png.Plate))

	if err != nil {
		log.Fatalf("Error while loading plate image: %v", err)
	}

	plateImage = ebiten.NewImageFromImage(plateStdImage)
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
	// Note that the IngredientType enum values MUST be kept ascending from 0 to N without any gaps for the
	// randomization to work correctly.
	randomIngredient := typeToIngredient[IngredientType(rand.IntN(len(typeToIngredient)))]

	// Spawn the Part in the middle lane.
	lane := laneCount / 2

	return newPart(randomIngredient, lane, buildScaleFactor)
}

func (p *Part) move(x, y int) {
	p.x = x
	p.y = y
}

func (p *Part) height() int {
	// Account for the scale factor when calculating the images' height.
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
	stack := make([]*Part, partCount)

	// Always start with a bottom bun and finsh with a top bun.
	stack[0] = newPart(typeToIngredient[IngBunBottom], lane, orderScaleFactor)
	stack[partCount-1] = newPart(typeToIngredient[IngBunTop], lane, orderScaleFactor)

	for index := 1; index < partCount-1; index++ {
		// Pick a random ingredient, avoid repeating a bottom or top bun.
		// Note that the IngredientType enum values MUST be kept ascending from 0 to N without any gaps for the
		// randomization to work correctly. Moreover, BunBottom and BunTop must be kept at the beginning of the enum.
		randomIngredient := typeToIngredient[IngredientType(rand.IntN(len(typeToIngredient)-2)+2)]

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

// ingredientTypes returns just the burgers ingredient types ordered from bottom to top.
func (b *Burger) ingredientTypes() []IngredientType {
	types := make([]IngredientType, len(b.stack))

	for index, part := range b.stack {
		types[index] = part.ingredient.typ
	}

	return types
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
