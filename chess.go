package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	err        error
	background *ebiten.Image
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	background, _, err = ebitenutil.NewImageFromFile("assets/table.jpeg", ebiten.FilterDefault)
	check(err)
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello, World!"); err != nil {
		log.Fatal(err)
	}
}
