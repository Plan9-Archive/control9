package control

import (
	"image"
	"strconv"
)

func ctlatoi(str string) int {
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return n
}

func ctlalignpt(r image.Rectangle, dx, dy, align int) image.Point {
	pt := image.ZP

	switch align % 3 {
	case 0: /* left */
		pt.X = r.Min.X
	case 1: /*center */
		pt.X = r.Min.X + (r.Dx()-dx)/2
	case 2: /*right*/
		pt.X = r.Max.X - dx
	}

	switch (align / 3) % 3 {
	case 0: /* top */
		pt.Y = r.Min.Y
	case 1: /* center */
		pt.Y = r.Min.Y + (r.Dy()-dy)/2
	case 2: /* bottom */
		pt.Y = r.Max.Y - dy
	}

	return pt
}
