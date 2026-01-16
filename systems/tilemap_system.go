package systems

import (
	"image"
	"image/png"
	"io"
	"os"

	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Tile struct {
	Coord math.Vector2
	Rect  utils.RectSpecs
	Color render.Color
}

type TileMap struct {
	Tiles      [][]*Tile
	instance   *lifecycle.GameObject
	can_render bool
}

func NewTileMap(map_size math.Vector2, tile_rect utils.RectSpecs, offset int) *TileMap {
	tiles := make([][]*Tile, 0)
	for y := range map_size.Y {
		row := make([]*Tile, 0)

		for x := range map_size.X {
			rect := tile_rect
			rect.PosX += x * (tile_rect.Width + offset)
			rect.PosY += y * (tile_rect.Height + offset)

			row = append(row, &Tile{
				Coord: math.Vector2{X: x, Y: y},
				Rect:  rect,
			})
		}

		tiles = append(tiles, row)
	}

	tilemap := &TileMap{
		Tiles:      tiles,
		can_render: false,
	}

	tilemap.instance = lifecycle.Register(&lifecycle.GameObject{
		Render: tilemap.render,
	})

	return tilemap
}

func (m *TileMap) DrawFromFile(path string) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(path)
	if err != nil {
		os.Exit(1)
		panic(err)
	}
	defer file.Close()

	pixels, err := getPixels(file)
	if err != nil {
		os.Exit(1)
		panic(err)
	}

	for i, row := range m.Tiles {
		for j, tile := range row {
			tile.Color = render.Color{
				R: pixels[i][j].R,
				G: pixels[i][j].G,
				B: pixels[i][j].B,
				A: pixels[i][j].A,
			}
		}
	}

	m.can_render = true
}

func (m *TileMap) DrawMap() {
	m.can_render = true
}

func (m *TileMap) render() {
	if !m.can_render {
		return
	}

	for _, row := range m.Tiles {
		for _, tile := range row {
			render.DrawRect(tile.Rect, tile.Color)
		}
	}
}

func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := range height {
		var row []Pixel
		for x := range width {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{
		uint8(r / 257),
		uint8(g / 257),
		uint8(b / 257),
		uint8(a / 257),
	}
}
