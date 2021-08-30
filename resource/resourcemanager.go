package resource

import (
	"errors"
	"fmt"
	"image"
	"log"

	"github.com/go-fonts/liberation/liberationsansregular"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jaeg/cool_game/config"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var Textures map[string]*ebiten.Image
var Fonts map[string]font.Face

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

func LoadFont(name string, path string) error {
	if Fonts == nil {
		log.Print("Initialize resource manager")
		Fonts = make(map[string]font.Face)
	}
	raw := liberationsansregular.TTF
	tt, err := opentype.Parse(raw)
	if err != nil {
		return err
	}

	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     config.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	Fonts[name] = fontFace
	return nil
}
