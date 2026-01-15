package values

import (
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
)

var SCREEN_SIZE math.Vector2 = math.Vector2{X: 1600, Y: 960}

var FONT_SPECS render.FontSpecs = render.FontSpecs{
	Name: "Pixelboy",
	Path: "assets/fonts/pixeboy-font/Pixeboy-z8XGD.ttf",
	Size: 32,
}
