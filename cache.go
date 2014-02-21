package control

import (
	"bitbucket.org/mischief/draw9"
	"bitbucket.org/mischief/draw9/color9"
	"image"
)

type imagecache map[string]*draw9.Image

var (
	coltab = map[string]color9.Color{
		"red":        color9.DRed,
		"blue":       color9.DBlue,
		"yellow":     color9.DYellow,
		"paleyellow": color9.DPaleyellow,
	}

	icache imagecache = map[string]*draw9.Image{}
)

func geti(name string) *draw9.Image {
	return icache[name]
}

// InitControlset must be called before any other control function to allocate colors.
func InitControlset(d *draw9.Display) error {
	var err error
	icache["opaque"] = d.Opaque
	icache["transparent"] = d.Transparent
	icache["white"] = d.White
	icache["black"] = d.Black
	for n, c := range coltab {
		icache[n], err = d.AllocImage(image.Rect(0, 0, 1, 1), d.ScreenImage.Pix, true, c)
		if err != nil {
			return err
		}
	}

	return err
}
