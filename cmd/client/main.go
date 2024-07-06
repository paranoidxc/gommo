package main

import (
	"context"
	"fmt"
	"gommo"
	"gommo/engine/asset"
	"gommo/engine/ecs"
	"gommo/engine/render"
	"gommo/engine/tilemap"
	"log"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"nhooyr.io/websocket"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// setup network
	url := "ws://localhost:8000"
	ctx := context.Background()
	c, resp, err := websocket.Dial(ctx, url, nil)
	check(err)

	log.Println("Connection Response:", resp)

	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)
	go func() {
		counter := byte(0)
		for {
			time.Sleep(1 * time.Second)
			n, err := conn.Write([]byte{counter})
			if err != nil {
				log.Println("Error Sending:", err)
				return
			}

			log.Println("Sent n Bytes:", n)
			counter++
		}
	}()

	// start pixel
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

	engine := ecs.NewEngine()

	load := asset.NewLoad(os.DirFS("./"))
	//load := asset.NewLoad(os.DirFS("./images"))
	//spritesheet, err := load.Spritesheet("packed.json")
	spritesheet, err := load.Spritesheet("spritesheet.json")
	check(err)

	// Create Tilemap
	seed := time.Now().UTC().UnixNano()

	mapSize := 1000
	tileSize := 16
	tmap := gommo.CreateTilemap(seed, mapSize, tileSize)

	grassTile, err := spritesheet.Get("grass0.png")
	check(err)
	dirtTile, err := spritesheet.Get("dirt0.png")
	check(err)
	waterTile, err := spritesheet.Get("water0.png")
	check(err)

	tmapRender := render.NewTilemapRender(spritesheet, map[tilemap.TileType]*pixel.Sprite{
		gommo.GrassTile: grassTile,
		gommo.DirtTile:  dirtTile,
		gommo.WaterTile: waterTile,
	})
	tmapRender.Batch(tmap)

	// create people
	spawnPoint := Transform{
		float64(tileSize * mapSize / 2),
		float64(tileSize * mapSize / 2)}

	//manSprite, err := load.Sprite("man.png")
	manSprite, err := spritesheet.Get("man_0.png")
	check(err)

	hatManSprite, err := spritesheet.Get("man_1.png")
	check(err)

	manId := engine.NewId()
	ecs.Write(engine, manId, Sprite{manSprite})
	ecs.Write(engine, manId, spawnPoint)
	ecs.Write(engine, manId, Keybinds{
		Up:    pixelgl.KeyUp,
		Down:  pixelgl.KeyDown,
		Left:  pixelgl.KeyLeft,
		Right: pixelgl.KeyRight,
	})

	hatManId := engine.NewId()
	ecs.Write(engine, hatManId, Sprite{hatManSprite})
	ecs.Write(engine, hatManId, spawnPoint)
	ecs.Write(engine, hatManId, Keybinds{
		Up:    pixelgl.KeyW,
		Down:  pixelgl.KeyS,
		Left:  pixelgl.KeyA,
		Right: pixelgl.KeyD,
	})

	camera := render.NewCamera(win, 0, 0)
	zoomSpeed := 0.1

	fmt.Println("Game Start")
	for !win.JustPressed(pixelgl.KeyEscape) {
		win.Clear(pixel.RGB(0, 0, 0))

		scroll := win.MouseScroll()
		if scroll.Y != 0 {
			camera.Zoom += zoomSpeed * scroll.Y
		}

		HandleInput(win, engine)

		transform := Transform{}
		ok := ecs.Read(engine, manId, &transform)
		if ok {
			camera.Position = pixel.V(transform.X, transform.Y)
		}
		camera.Update()

		win.SetMatrix(camera.Mat())
		// render first because tile behide the people
		tmapRender.Draw(win)

		DrawSprites(win, engine)

		win.SetMatrix(pixel.IM)

		win.Update()
	}
	fmt.Println("Game Quit")
}

type Keybinds struct {
	Up, Down, Left, Right pixelgl.Button
}

func (t *Keybinds) ComponentSet(val interface{}) {
	*t = val.(Keybinds)
}

type Sprite struct {
	*pixel.Sprite
}

func (t *Sprite) ComponentSet(val interface{}) {
	*t = val.(Sprite)
}

type Transform struct {
	X, Y float64
}

func (t *Transform) ComponentSet(val interface{}) {
	*t = val.(Transform)
}

func DrawSprites(win *pixelgl.Window, engine *ecs.Engine) {
	ecs.Each(engine, Sprite{}, func(id ecs.Id, a interface{}) {
		sprite := a.(Sprite)

		transform := Transform{}
		ok := ecs.Read(engine, id, &transform)
		if !ok {
			return
		}

		pos := pixel.V(transform.X, transform.Y)
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pos))
	})
}

/*
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
*/

func HandleInput(win *pixelgl.Window, engine *ecs.Engine) {
	ecs.Each(engine, Keybinds{}, func(id ecs.Id, a interface{}) {
		keybinds := a.(Keybinds)

		transform := Transform{}
		ok := ecs.Read(engine, id, &transform)
		if !ok {
			return
		}

		if win.Pressed(keybinds.Left) {
			transform.X -= 2.0
		}
		if win.Pressed(keybinds.Right) {
			transform.X += 2.0
		}
		if win.Pressed(keybinds.Up) {
			transform.Y += 2.0
		}
		if win.Pressed(keybinds.Down) {
			transform.Y -= 2.0
		}

		ecs.Write(engine, id, transform)
	})
}
