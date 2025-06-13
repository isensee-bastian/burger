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
		path := fmt.Sprintf("%s/%s.png", imageBasePath, ingredient.name)
		image, _, err := ebitenutil.NewImageFromFile(path)

		if err != nil {
			log.Fatalf("Error while loading ingredient image %s: %v", ingredient.name, err)
		}

		allIngredients[index].image = image
	}
}

// Part is a specific instances of an ingredient with a specific position on the screen. Note that there can be multiple
// Parts of the same Ingredient, but at different locations on the screen, e.g. multiple pickle instances.
type Part struct {
	ingredient Ingredient
	lane       int
	x          int
	y          int
}

func newPart(ingredient Ingredient) *Part {
	// Spawn the Part in the middle lane.
	lane := laneCount / 2

	return &Part{
		ingredient: ingredient,
		lane:       lane,
		x:          lane * laneWidth,
		y:          0,
	}
}

func newRandomPart() *Part {
	if len(allIngredients) == 0 {
		log.Fatal("Ingredients not initialized yet, cannot create random part")
	}
	randomIngredient := allIngredients[rand.IntN(len(allIngredients))]

	return newPart(randomIngredient)
}

func (p *Part) move(x, y int) {
	p.x = x
	p.y = y
}

func (p *Part) height() int {
	// Account for the scale factor when calculating the images height.
	return int(float64(p.ingredient.image.Bounds().Size().Y) * scaleFactor)
}

func (p *Part) draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleFactor, scaleFactor)
	opts.GeoM.Translate(float64(p.x), float64(p.y))

	screen.DrawImage(p.ingredient.image, opts)
}

// Burger is a stack of Part instances, i.e. specific ingredients at a certain location ordered from bottom to top.
type Burger struct {
	stack []*Part
}

func newBurger() *Burger {
	return &Burger{stack: make([]*Part, 0)}
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
	for _, part := range b.stack {
		part.draw(screen)
	}
}
