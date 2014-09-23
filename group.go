package control

import (
	"bitbucket.org/mischief/draw9"
	"fmt"
	"image"
	"log"
	"strings"
)

type Group struct {
	*BaseControl
	lastbut, border int
	mansize         bool
	separation      int
	selected        int
	lastkid         int
	bordercolor     *draw9.Image
	image           *draw9.Image
	kids            []Control
	separators      []image.Rectangle
}

func (g *Group) Ctl(e string) {
	var r image.Rectangle
	f := strings.Fields(e)

	switch f[0] {
	case "add":
		for i := 1; i < len(f); i++ {
			if c := g.set.Called(f[i]); c != nil {
				g.add(c)
			} else {
				log.Printf("no such control %q", f[i])
			}
		}
	case "border":
		bsz := ctlatoi(f[1])
		if bsz < 0 {
			log.Printf("bad border size %d", bsz)
		} else {
			g.border = ctlatoi(f[1])
		}
	case "bordercolor":
	case "focus":
	case "hide":
		for _, c := range g.kids {
			c.Ctl("hide")
		}
		g.hidden = true
	case "image":
	case "rect":
		if len(f) < 5 {
			log.Printf("%s: bad rectangle: %s", g.name, e)
		} else {
			r := image.Rect(0, 0, 1, 1)
			r.Min.X = ctlatoi(f[1])
			r.Min.Y = ctlatoi(f[2])
			r.Max.X = ctlatoi(f[3])
			r.Max.Y = ctlatoi(f[4])

			if r.Dx() <= 0 || r.Dy() <= 0 {
				log.Printf("%s: bad rectangle: %s", g.name, e)
				break
			}

			g.rect = r
			r = r.Inset(g.border)
			if len(g.kids) == 0 {
				return
			}
			switch g.typ {
			case CtlRow:
				g.rowresize(r)
			}
		}
	case "remove":
	case "reveal":
	case "show":
		if g.hidden {
			break
		}
		g.screen.Border(g.rect, g.border, g.bordercolor, g.bordercolor.R.Min)
		r := g.rect.Inset(g.border)
		g.screen.Draw(r, g.image, nil, g.image.R.Min)
		for i := 0; i < len(g.kids); i++ {
			g.kids[i].Ctl("show")
		}
		g.image.Display.Flush()
	case "size":
		switch len(f) {
		case 1:
			g.mansize = false
			g.SetSize()
		case 5:
			r.Max.X = ctlatoi(f[3])
			r.Max.Y = ctlatoi(f[4])
			fallthrough
		case 3:
			r.Min.X = ctlatoi(f[1])
			r.Min.Y = ctlatoi(f[2])
			if r.Min.X <= 0 || r.Min.Y <= 0 || r.Max.X <= 0 || r.Max.Y <= 0 || r.Max.X < r.Min.X || r.Max.Y < r.Min.Y {
				log.Printf("%s: bad rectangle: %s", g.name, e)
				return
			}
			g.size = r
			g.mansize = true
		}
	case "separation":
	default:
		log.Printf("%q: unrecognized message %q", g.name, f[0])
	}
}

// interface functions

func (g *Group) Mouse(ms draw9.Mouse) {
	log.Printf("%s mouse %s", g.name, ms)
	/* TODO CtlStack */
	lastkid := -1
	for i, k := range g.kids {
		if ((!ms.Any() || g.lastbut == 0) && ms.Pt().In(k.Rect())) || ((ms.Any() || g.lastbut != 0) && g.lastkid == i) {
			k.Mouse(ms)
			lastkid = i
		}
	}
	g.lastkid = lastkid
	g.lastbut = ms.Buttons()
}

func (g *Group) Key(r rune) {}
func (g *Group) Exit()      {}
func (g *Group) SetSize() {
	if g.mansize {
		return
	}

	r := image.Rect(1, 1, 1, 1)
	for i, k := range g.kids {
		k.SetSize()
		sz := k.Size()
		if sz.Min.X == 0 || sz.Min.Y == 0 || sz.Max.X == 0 || sz.Max.Y == 0 {
			log.Printf("%s: invalid size: %s %s", g.name, k.Name(), sz)
			return
		}
		switch g.typ {
		case CtlRow:
			if i > 0 {
				r.Min.X += sz.Min.X + g.border
				r.Max.X += sz.Max.X + g.border
			} else {
				r.Min.X = sz.Min.X
				r.Max.X = sz.Max.X
			}
			if r.Min.Y < sz.Min.Y {
				r.Min.X = sz.Min.Y
			}
			if r.Max.Y < sz.Max.Y {
				r.Max.Y = sz.Max.Y
			}
		}
	}
	g.size = r.Add(image.Pt(g.border, g.border))
	log.Printf("%s size %s", g.name, g.size)
}
func (g *Group) Activate() {}

// private functions

func (g *Group) init() {
	g.bordercolor = geti("black")
	g.image = geti("white")
	g.border = 0
	g.mansize = false
	g.separation = 0
	g.selected = -1
	g.lastkid = -1
}

func (g *Group) add(c Control) {
	g.kids = append(g.kids, c)
}

// resizing functions
func (g *Group) rowresize(r image.Rectangle) {
	var j int
	var rr image.Rectangle

	x := r.Dx()
	y := r.Dy()

	if x < g.size.Min.X {
		log.Printf("rowresize %s: too narrow: need %d, have %d", g.name, g.size.Min.X, x)
		r.Max.X = r.Min.X + g.size.Min.X
		x = r.Dx()
	}

	if y < g.size.Min.Y {
		log.Printf("rowresize %s: too short: need %d, have %d", g.name, g.size.Min.Y, y)
		r.Max.Y = r.Min.Y + g.size.Min.Y
	}

	d := make([]int, len(g.kids))
	p := make([]int, len(g.kids))
	for i, k := range g.kids {
		ksz := k.Size()
		d[i] = ksz.Min.X
		x -= d[i]
		p[i] = ksz.Max.Y - ksz.Min.X
	}
	x -= (len(g.kids) - 1) * g.separation
	if x < 0 {
		log.Printf("rowresize %s: x == %d", g.name, x)
		x = 0
	}

	if x >= g.size.Max.X-g.size.Min.X {
		log.Printf("rowresize %s: max: %d > %d - %d", g.name, x, g.size.Max.X, g.size.Min.X)
		for i, _ := range g.kids {
			d[i] += p[i]
		}
		x -= g.size.Max.X - g.size.Min.X
	} else {
		log.Printf("rowresize %s: divvie up: %d < %d - %d", g.name, x, g.size.Max.X, g.size.Min.X)
		// rects cant be max width
		j = x
		for i, _ := range g.kids {
			t := p[i] * x / (g.size.Max.X - g.size.Min.X)
			d[i] += t
			j -= t
		}
		d[0] += j
		x = 0
	}
	g.separators = make([]image.Rectangle, len(g.kids)-1)
	rr = r
	for i, k := range g.kids {
		if i < len(g.kids)-1 {
			g.separators[i].Min.Y = r.Min.Y
			g.separators[i].Max.Y = r.Max.Y
		}
		t := x / (len(g.kids) - i)
		x -= t
		j += t / 2

		rr.Min.X = r.Min.X + j
		if i > 0 {
			g.separators[i-1].Max.X = rr.Min.X
		}
		j += d[i]
		rr.Max.X = r.Min.X + j
		if i < len(g.kids)-1 {
			g.separators[i].Min.X = rr.Max.X
		}
		j += g.separation + t - t/2
		k.Ctl(fmt.Sprintf("rect %d %d %d %d", rr.Min.X, rr.Min.Y, rr.Max.X, rr.Max.Y))

	}
}

// ctor

// Row creates a new Group that is visually represented as a row of controls.
func (cs *Controlset) Row(name string) Control {
	row := &Group{
		BaseControl: cs.newctl(name, CtlRow),
	}

	row.init()

	cs.addctl(row)

	return row
}
