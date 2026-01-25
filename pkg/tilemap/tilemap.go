package tilemap

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/mikabrytu/gomes-engine/lifecycle"
	"github.com/mikabrytu/gomes-engine/math"
	"github.com/mikabrytu/gomes-engine/render"
	"github.com/mikabrytu/gomes-engine/utils"
)

type Channel int

const (
	R Channel = iota
	G
	B
	A
)

type Pixel struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type Tile struct {
	Coord      math.Vector2
	Rect       utils.RectSpecs
	Sprite     *render.Sprite
	Color      render.Color
	HasEnemy   bool
	HasItem    bool
	IsWalkable bool
	Enabled    bool
}

type TileRules struct {
	Chan       Channel
	ChanValue  uint8
	SpritePath string
	Color      render.Color
	Walkable   bool
}

type TileMap struct {
	Tiles       [][]*Tile
	instance    *lifecycle.GameObject
	render_rect bool
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
				Coord:    math.Vector2{X: x, Y: y},
				Rect:     rect,
				HasEnemy: false,
				HasItem:  false,
				Enabled:  false,
			})
		}

		tiles = append(tiles, row)
	}

	tilemap := &TileMap{
		Tiles:       tiles,
		render_rect: false,
	}

	tilemap.instance = lifecycle.Register(&lifecycle.GameObject{
		Render:  tilemap.render,
		Destroy: tilemap.destroy,
	})

	return tilemap
}

func (m *TileMap) Enable() {
	for _, row := range m.Tiles {
		for _, tile := range row {
			if tile.Sprite != nil && tile.Enabled {
				tile.Sprite.Enable()
			}
		}
	}

	lifecycle.Enable(m.instance)
}

func (m *TileMap) Disable() {
	for _, row := range m.Tiles {
		for _, tile := range row {
			if tile.Sprite != nil {
				tile.Sprite.Disable()
			}
		}
	}

	lifecycle.Disable(m.instance)
}

func (m *TileMap) DrawMapFile(file string) {
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

	m.render_rect = true
}

func (m TileMap) DrawMapAssetsFromFile(rules []TileRules, file string) {
	pixels, err := getPixels(file)
	if err != nil {
		panic(err)
	}

	for i, row := range m.Tiles {
		for j, tile := range row {
			for _, r := range rules {
				switch r.Chan {
				case R:
					if pixels[i][j].R == r.ChanValue {
						tile.HasEnemy = true
						tile.Enabled = false
						tile.Color = r.Color
					}
				case G:
					if pixels[i][j].G == r.ChanValue {
						tile.HasItem = true
						tile.Enabled = false
						tile.Color = r.Color
					}
				case B:
					if pixels[i][j].B == r.ChanValue {
						tile.IsWalkable = r.Walkable

						if r.SpritePath == "" {
							continue
						}

						sprite := render.NewSprite(
							fmt.Sprintf("tile-%v-%v-tile", i, j),
							r.SpritePath,
							tile.Rect,
							render.White,
						)
						sprite.Init()
						sprite.Disable()

						tile.Sprite = sprite
					}
				}
			}
		}
	}
}

func (m *TileMap) DrawMap() {
	m.render_rect = true
}

func (m *TileMap) render() {
	for _, row := range m.Tiles {
		for _, tile := range row {
			rect := utils.RectSpecs{
				PosX:   tile.Rect.PosX + (tile.Rect.Width / 4),
				PosY:   tile.Rect.PosY + (tile.Rect.Height / 4),
				Width:  tile.Rect.Width / 2,
				Height: tile.Rect.Height / 2,
			}

			if m.render_rect {
				render.DrawRect(tile.Rect, tile.Color)
			}

			if tile.Enabled {
				if tile.HasEnemy {
					render.DrawRect(rect, tile.Color)
				}

				if tile.HasItem {
					render.DrawRect(rect, tile.Color)
				}
			}
		}
	}
}

func (m *TileMap) destroy() {
	for _, row := range m.Tiles {
		for _, tile := range row {
			tile.Sprite.Clear()
		}
	}
}

func getPixels(path string) ([][]Pixel, error) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open(path)
	if err != nil {
		os.Exit(1)
		panic(err)
	}
	defer file.Close()

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

	//fmt.Println(pixels)

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
