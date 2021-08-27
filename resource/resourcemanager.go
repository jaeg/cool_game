package resource

import (
	"errors"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Textures map[string]*ebiten.Image

func LoadImageAsTexture(name string, path string) error {
	if Textures == nil {
		log.Print("Initialize resource manager")
		Textures = make(map[string]*ebiten.Image)
	}
	img, err := LoadImage(path)
	if err != nil {
		return err
	}

	Textures[name] = img
	return nil
}

func LoadImage(path string) (*ebiten.Image, error) {
	imgFile, err := ebitenutil.OpenFile(path)
	if err != nil {
		fmt.Println("Error opening tileset " + path)
		return nil, errors.New("error opening tileset " + path)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
