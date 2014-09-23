package control

import (
	"bitbucket.org/mischief/draw9"
	"fmt"
	"image"
	"log"
	"strings"
)

type Button struct {
	*BaseControl
	image, mask                  *draw9.Image
	light, pale, bordercolor     *draw9.Image
	pressed, lastbut, lastshow   int
	border, align, off, prepress int
}

// Control interface functions

func (b *Button) Ctl(e string) {
	log.Printf("%s ctl %s", b.name, e)
	f := strings.Fields(e)
	switch f[0] {
	default:
		log.Printf("%s: unrecognized message %s", b.name, e)
	case "align":
	case "border":
	case "bordercolor":
	case "focus":
	case "format":
	case "hide":
	case "image":
	case "light":
	case "mask":
	case "pale":
	case "rect":
	case "reveal":
		b.hidden = false
		fallthrough
	case "show":
		b.show()
	case "size":
	case "value":
		arg := ctlatoi(f[1])
		if arg != b.pressed {
			b.pressed ^= 1
			b.show()
		}
	}
}
func (b *Button) Mouse(m draw9.Mouse) {
	log.Printf("%s mouse %s", b.name, m)
	if m.Any() {
		/* some button is down */
		if m.Pt().In(b.rect) {
			b.off = 0
			b.show()
		} else {
			if b.off == 0 {
				b.off = 1
				b.show()
			}
		}
	}
	if m.Buttons() & 7 != b.lastbut {
		/* button change */
		if m.Buttons() & 7 > 0 {
			b.prepress = b.pressed
			if b.pressed > 0 {
				b.pressed = 0
			} else {
				b.pressed = m.Buttons() & 7
			}
			b.show()
		} else {
			if m.Pt().In(b.rect) && b.event != nil {
				b.event <- fmt.Sprintf("%s: value %d", b.name, b.pressed)
			} else {
				b.off = 0
				b.pressed = b.prepress
				b.show()
			}
		}
	}

	b.lastbut = m.Buttons() & 7
}

func (b *Button) Key(c rune) {}
func (b *Button) Exit()      {}
func (b *Button) SetSize()   {}
func (b *Button) Activate()  {}

// Button-specific

// private junk

// draw button
func (b *Button) show() {
	if b.hidden {
		return
	}

	r := b.rect

	if b.border > 0 {
		b.screen.Border(r, b.border, b.bordercolor, image.ZP)
		r = r.Inset(b.border)
	}

	b.screen.Draw(r, b.image, nil, b.image.R.Min)
	if b.off > 0 {
		b.screen.Draw(r, b.pale, b.mask, b.mask.R.Min)
	} else if b.pressed > 0 {
		b.screen.Draw(r, b.light, b.mask, b.mask.R.Min)
	}

	b.lastshow = b.pressed
	b.image.Display.Flush()
}

// ctor

// Button creates a new Button Control with the given name.
func (cs *Controlset) Button(name string) Control {
	b := &Button{
		BaseControl: cs.newctl(name, CtlButton),
		image:       geti("white"),
		mask:        geti("transparent"),
		light:       geti("yellow"),
		pale:        geti("paleyellow"),
		bordercolor: geti("black"),
		border:      2,
		align:       0,
		off:         0,
		prepress:    0,
	}

	cs.addctl(b)

	return b
}
