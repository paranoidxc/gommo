package render

import (
	"gommo/engine/asset"
	"gommo/engine/tilemap"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type TilemapRender struct {
	spritesheet  *asset.Spritesheet
	batch        *pixel.Batch
	tileToSprite map[tilemap.TileType]*pixel.Sprite
}

func NewTilemapRender(spritesheet *asset.Spritesheet, tileToSprite map[tilemap.TileType]*pixel.Sprite) *TilemapRender {
	return &TilemapRender{
		spritesheet:  spritesheet,
		batch:        pixel.NewBatch(&pixel.TrianglesData{}, spritesheet.Picture()),
		tileToSprite: tileToSprite,
	}
}

func (r *TilemapRender) Clear() {
	r.batch.Clear()
}

func (r *TilemapRender) Batch(t *tilemap.Tilemap) {
	for x := 0; x < t.Width(); x++ {
		for y := 0; y < t.Height(); y++ {
			tile, ok := t.Get(x, y)
			if !ok {
				continue
			}
			pos := pixel.V(float64(x*t.TileSize),
				float64(y*t.TileSize))

			sprite, ok := r.tileToSprite[tile.Type]
			if !ok {
				panic("Unable to find TileType")
			}
			mat := pixel.IM.Moved(pos)
			sprite.Draw(r.batch, mat)
		}
	}

}

func (t *TilemapRender) Draw(win *pixelgl.Window) {
	t.batch.Draw(win)
}
