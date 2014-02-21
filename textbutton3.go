package control

import (
	"bitbucket.org/mischief/draw9"
	"fmt"
	"image"
	"log"
	"strings"
)

type TextButton3 struct {
	*BaseControl
	font                        *draw9.Font
	image, mask                 *draw9.Image
	light, bordercolor    *draw9.Image
	textcolor, pressedtextcolor *draw9.Image
	lastbut                     draw9.Mouse
	pressed, lastshow           int
	lines                       []string
	border, align               int
	left, middle, right         int
	toggle, gettextflg          int
}

// Control interface functions

func (b *TextButton3) Ctl(e string) {
	log.Printf("%s ctl %s", b.name, e)
	f := strings.Fields(e)
	switch f[0] {
	default:
		log.Printf("%s: unrecognized message %s", b.name, e)
	case "align":
		b.align = ctlatoi(f[1])
		b.lastshow = -1
	case "border":
		b.border = ctlatoi(f[1])
		b.lastshow = -1
	case "bordercolor":
		b.bordercolor = geti(f[1])
		b.lastshow = -1
	case "focus":
	case "font":
		/* fixme */
		b.lastshow = -1
	case "format":
		/* fixme */
	case "hide":
		b.hidden = true
	case "image":
		b.image = geti(f[1])
		b.lastshow = -1
	case "light":
		b.light = geti(f[1])
		b.lastshow = -1
	case "mask":
		b.mask = geti(f[1])
		b.lastshow = -1
	case "pressedcolor":
		b.pressedtextcolor = geti(f[1])
		b.lastshow = -1
	case "rect":
		if len(f) < 5 {
			log.Printf("%s: bad rectangle: %s", b.name, e)
			return
		}
		r := image.Rect(ctlatoi(f[1]), ctlatoi(f[2]), ctlatoi(f[3]), ctlatoi(f[4]))
		if r.Dx() <= 0 || r.Dy() <= 0 {
			log.Printf("%s: bad rectangle: %s", b.name, e)
			return
		}
		b.rect = r
		b.lastshow = -1
	case "reveal":
		b.hidden = false
		b.lastshow = -1
		b.show()
	case "show":
		b.lastshow = -1
		b.show()
	case "size":
		if len(f) != 3 && len(f) != 5 {
			log.Printf("%s: bad rectangle: %s", b.name, e)
			return
		}
		r := image.Rect(ctlatoi(f[1]), ctlatoi(f[2]), 0x7fffffff, 0x7fffffff)
		if len(f) == 5 {
			r.Max.X = ctlatoi(f[3])
			r.Max.Y = ctlatoi(f[4])
		}
		if r.Min.X <= 0 || r.Min.Y <= 0 || r.Max.X <= 0 || r.Max.Y <= 0 || r.Max.X < r.Min.X || r.Max.Y < r.Min.Y {
			log.Printf("%s: bad rectangle: %s", b.name, e)
			return
		}
		b.size = r
	case "text":
		b.lines = strings.Split(strings.Join(f[1:], ""), "\n")
		b.lastshow = -1
		b.show()
	case "textcolor":
		b.textcolor = geti(f[1])
		b.lastshow = -1
	case "value":
		arg := ctlatoi(f[1])
		if arg != b.pressed {
			b.pressed ^= 1
			b.show()
		}
	}
}
func (b *TextButton3) Mouse(m draw9.Mouse) {
	log.Printf("%s mouse %s", b.name, m)

	if b.left == 1 {
		if m.Mb1() && !b.lastbut.Mb1() {
			b.pressed ^= 1
			b.show()
			b.lastbut = m
		} else if !m.Mb1() && b.lastbut.Mb1() {
			if b.gettextflg == 0 {
				b.event <- fmt.Sprintf(b.format, b.name, b.pressed, m.Pt().X, m.Pt().Y)
			} else {
				b.event <- fmt.Sprintf("%s: value %s", b.name, b.lines[0])
			}
			b.pressed ^= 1
			b.show()
			b.lastbut = m
		}
	}

	if b.middle == 1 {
		if m.Mb2() && !b.lastbut.Mb2() {
			b.pressed ^= 2
			b.show()
			b.lastbut = m
		} else if !m.Mb2() && b.lastbut.Mb2() {
			if b.gettextflg == 0 {
				b.event <- fmt.Sprintf(b.format, b.name, b.pressed, m.Pt().X, m.Pt().Y)
			} else {
				b.event <- fmt.Sprintf("%s: value %s", b.name, b.lines[0])
			}
			b.pressed ^= 2
			b.show()
			b.lastbut = m
		}
	}
	if b.right == 1 {
		if m.Mb3() && !b.lastbut.Mb3() {
			b.pressed ^= 4
			b.show()
			b.lastbut = m
		} else if !m.Mb3() && b.lastbut.Mb3() {
			if b.gettextflg == 0 {
				b.event <- fmt.Sprintf(b.format, b.name, b.pressed, m.Pt().X, m.Pt().Y)
			} else {
				b.event <- fmt.Sprintf("%s: value %s", b.name, b.lines[0])
			}
			b.pressed ^= 4
			b.show()
			b.lastbut = m
		}
	}
}

func (b *TextButton3) Key(c rune) {}
func (b *TextButton3) Exit()      {}
func (b *TextButton3) SetSize()   {}
func (b *TextButton3) Activate()  {}

// TextButton3-specific

// private junk

// draw button
func (b *TextButton3) show() {
	if b.hidden || b.lastshow == b.pressed {
		return
	}

	r := b.rect

	b.screen.Draw(r, b.image, nil, b.image.R.Min)
	if b.pressed > 0 || b.toggle > 0 {
		b.screen.Draw(r, b.light, b.mask, b.mask.R.Min)
	}

	if b.border > 0 {
		b.screen.Border(r, b.border, b.bordercolor, image.ZP)
		r = r.Inset(b.border)
	}

	dx := 0
	f := b.font

	for _, s := range b.lines {
		w := f.StringWidth(s)
		if dx < w {
			dx = w
		}
	}

	dy := len(b.lines) * f.Height
	clipr := b.rect.Inset(b.border)
	pt := ctlalignpt(clipr, dx, dy, b.align)
	im := b.textcolor
	if b.pressed > 0 {
		im = b.pressedtextcolor
	}

	for i := 0; i < len(b.lines); i++ {
		r.Min = pt
		r.Max.X = pt.X + dx
		r.Max.Y = pt.Y + f.Height
		q := ctlalignpt(r, f.StringWidth(b.lines[i]), f.Height, b.align%3)
		b.screen.String(q, im, image.ZP, f, b.lines[i])
		pt.Y += f.Height
	}

	b.lastshow = b.pressed
	b.image.Display.Flush()
}

// ctor

// TextButton3 creates a new TextButton3 Control with the given name.
func (cs *Controlset) TextButton3(name string) Control {
	b := &TextButton3{
		BaseControl:      cs.newctl(name, CtlTextButton3),
		image:            geti("white"),
		mask:             geti("transparent"),
		light:            geti("yellow"),
		bordercolor:      geti("black"),
		textcolor:        geti("black"),
		pressedtextcolor: geti("black"),
		font:             cs.screen.Display.DefaultFont,
		border:           2,
		align:            4,
		left:             1,
		middle:           1,
		right:            1,
	}

	b.format = "%s: value %d %d %d"

	cs.addctl(b)

	return b
}
