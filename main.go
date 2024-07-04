package main

import (
	"fmt"
	"gommo/engine/asset"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

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
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)

	load := asset.NewLoad(os.DirFS("./"))
	//load := asset.NewLoad(os.DirFS("./images"))
	spritesheet, err := load.Spritesheet("packed.json")
	if err != nil {
		panic(err)
	}

	//manSprite, err := load.Sprite("man.png")
	manSprite, err := spritesheet.Get("man.png")
	if err != nil {
		panic(err)
	}
	manPosition := win.Bounds().Center()

	fmt.Println("start")
	for !win.JustPressed(pixelgl.KeyEscape) {
		win.Clear(pixel.RGB(0, 0, 0))

		if win.Pressed(pixelgl.KeyLeft) {
			manPosition.X -= 2.0
		}
		if win.Pressed(pixelgl.KeyRight) {
			manPosition.X += 2.0
		}
		if win.Pressed(pixelgl.KeyUp) {
			manPosition.Y += 2.0
		}
		if win.Pressed(pixelgl.KeyDown) {
			manPosition.Y -= 2.0
		}

		manSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(manPosition))
		win.Update()

	}
}
