package main

import (
	"fmt"
	"gommo/engine/asset"
	"gommo/engine/render"
	"gommo/engine/tilemap"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pixelgl.Run(runGame)
}

func runGame() {
	cfg := pixelgl.WindowConfig{
		Title:     "MMO",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	check(err)

	win.SetSmooth(false)

	load := asset.NewLoad(os.DirFS("./"))
	//load := asset.NewLoad(os.DirFS("./images"))
	//spritesheet, err := load.Spritesheet("packed.json")
	spritesheet, err := load.Spritesheet("spritesheet.json")
	check(err)

	// Create Tilemap
	grassSprite, err := spritesheet.Get("grass0.png")
	check(err)
	tileSize := 16
	mapSize := 100
	tiles := make([][]tilemap.Tile, mapSize, mapSize)
	for x := range tiles {
		tiles[x] = make([]tilemap.Tile, mapSize, mapSize)
		for y := range tiles[x] {
			tiles[x][y] = tilemap.Tile{
				Type:   0,
				Sprite: grassSprite,
			}
		}
	}
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet.Picture())
	tmap := tilemap.New(tiles, batch, tileSize)
	tmap.Rebatch()

	// create people
	spawnPoint := pixel.V(
		float64(tileSize*mapSize/2),
		float64(tileSize*mapSize/2))

	//manSprite, err := load.Sprite("man.png")
	manSprite, err := spritesheet.Get("man_0.png")
	check(err)

	hatManSprite, err := spritesheet.Get("man_1.png")
	check(err)

	people := make([]Person, 0)
	people = append(people, NewPerson(manSprite, spawnPoint, Keybinds{
		Up:    pixelgl.KeyUp,
		Down:  pixelgl.KeyDown,
		Left:  pixelgl.KeyLeft,
		Right: pixelgl.KeyRight,
	}))
	people = append(people, NewPerson(hatManSprite, spawnPoint, Keybinds{
		Up:    pixelgl.KeyW,
		Down:  pixelgl.KeyS,
		Left:  pixelgl.KeyA,
		Right: pixelgl.KeyD,
	}))

	camera := render.NewCamera(win, 0, 0)
	zoomSpeed := 0.1

	fmt.Println("Game Start")
	for !win.JustPressed(pixelgl.KeyEscape) {
		win.Clear(pixel.RGB(0, 0, 0))

		scroll := win.MouseScroll()
		if scroll.Y != 0 {
			camera.Zoom += zoomSpeed * scroll.Y
		}

		// 不能使用 for kv range
		for i := range people {
			people[i].HandleInput(win)
		}
		// Collison Detection

		camera.Position = people[0].Position
		camera.Update()

		win.SetMatrix(camera.Mat())
		// render first because tile behide the people
		tmap.Draw(win)
		for i := range people {
			people[i].Draw(win)
		}
		win.SetMatrix(pixel.IM)

		win.Update()
	}
	fmt.Println("Game Quit")
}

const (
	GrassTile tilemap.TileType = iota
	DirtTile
	WaterTile
)

func GetTile(ss *asset.Spritesheet, t tilemap.TileType) tilemap.Tile {
	name := ""
	switch t {
	case GrassTile:
		name = "grass0.png"
	case DirtTile:
		name = "grass0.png"
	case WaterTile:
		name = "water.png"
	default:
		panic("Unknow TileType")
	}

	sprite, err := ss.Get(spriteName)
	check(err)
}

type Keybinds struct {
	Up, Down, Left, Right pixelgl.Button
}

type Person struct {
	Sprite   *pixel.Sprite
	Position pixel.Vec
	Keybinds Keybinds
}

func NewPerson(sprite *pixel.Sprite,
	position pixel.Vec,
	keybinds Keybinds,
) Person {
	return Person{
		Sprite:   sprite,
		Position: position,
		Keybinds: keybinds,
	}
}

func (p *Person) Draw(win *pixelgl.Window) {
	p.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(p.Position))
}

func (p *Person) HandleInput(win *pixelgl.Window) {
	if win.Pressed(p.Keybinds.Left) {
		p.Position.X -= 2.0
	}
	if win.Pressed(p.Keybinds.Right) {
		p.Position.X += 2.0
	}
	if win.Pressed(p.Keybinds.Up) {
		p.Position.Y += 2.0
	}
	if win.Pressed(p.Keybinds.Down) {
		p.Position.Y -= 2.0
	}
}
