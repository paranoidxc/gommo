package tilemap

type TileType uint8 // 256

type Tile struct {
	Type TileType
}

type Tilemap struct {
	TileSize int
	tiles    [][]Tile
}

func New(tiles [][]Tile, tileSize int) *Tilemap {
	return &Tilemap{
		TileSize: tileSize,
		tiles:    tiles,
	}
}

func (t *Tilemap) Width() int {
	return len(t.tiles)
}

func (t *Tilemap) Height() int {
	// TODO assumes the tile map is a square
	// and is large than size 0
	return len(t.tiles[0])
}

func (t *Tilemap) Get(x, y int) (Tile, bool) {
	if x < 0 || x >= len(t.tiles) || y < 0 || y >= len(t.tiles[0]) {
		return Tile{}, false
	}
	return t.tiles[x][y], true
}
