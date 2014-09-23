package control

import (
	"bitbucket.org/mischief/draw9"
	"bufio"
	"fmt"
	"image"
	"strings"
)

type Controlset struct {
	controls    []Control
	screen      *draw9.Image
	active      []Control
	focus       Control
	ctl         chan string
	data        chan string
	kbdc        chan rune
	mousec      chan draw9.Mouse
	resizec     chan int
	resizeexitc chan int
	OnResize    func(cs *Controlset)
	csexitc     chan int
	kbd         *draw9.Consctl
	ms          *draw9.Mousectl
	clicktotype bool
}

func NewControlset(screen *draw9.Image, kbd *draw9.Consctl, ms *draw9.Mousectl) (*Controlset, error) {
	if screen == nil {
		return nil, fmt.Errorf("nil screen")
	}
	if kbd == nil {
	}
	if ms == nil {
	}

	cs := &Controlset{
		screen:      screen,
		ctl:         make(chan string, 64),
		data:        make(chan string, 64),
		kbdc:        kbd.C,
		mousec:      ms.C,
		resizec:     ms.Resize,
		resizeexitc: make(chan int),
		csexitc:     make(chan int),
		kbd:         kbd,
		ms:          ms,
	}

	go cs.resize()
	go cs.cs()

	return cs, nil
}

// Called looks up a child control by name.
func (cs *Controlset) Called(name string) Control {
	for _, c := range cs.controls {
		if c.Name() == name {
			return c
		}
	}

	return nil
}

func (cs *Controlset) newctl(name string, typ ControlType) *BaseControl {
	bc := &BaseControl{
		name: name,
		size: image.Rect(1, 1, 10000, 10000),
		//rect: image.Rect(1, 1, 100, 100),
		//		event:  make(chan string, 64),
		data:   make(chan string, 0),
		typ:    typ,
		hidden: false,
		set:    cs,
		screen: cs.screen,
	}

	return bc
}

func (cs *Controlset) addctl(c Control) {
	cs.controls = append(cs.controls, c)
	cs.active = append(cs.active, c)
}

func (cs *Controlset) Ctl(s string) {
	cs.ctl <- s
}

func (cs *Controlset) cmd(s string) {
	f := strings.Fields(s)
	for _, c := range cs.controls {
		if c.Name() == f[0] {
			c.Ctl(strings.Join(f[1:], " "))
		}
	}
}

func (cs *Controlset) onctl(s string) {
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		cs.cmd(sc.Text())
	}
}

func (cs *Controlset) Resize() {
	if cs.OnResize != nil {
		cs.OnResize(cs)
	}
}

func (cs *Controlset) resize() {
	for {
		select {
		case <-cs.resizec:
			cs.Resize()
		case <-cs.resizeexitc:
			return
		}
	}
}

func (cs *Controlset) cs() {
	prevbut := 0

	for {
		select {
		case r := <-cs.kbdc:
			if cs.focus != nil {
				cs.focus.Key(r)
			}
		case ms := <-cs.mousec:
			if prevbut > 0 {
				goto Send
			}
			if cs.focus != nil && cs.focus.Hidden() == false && ms.Pt().In(cs.focus.Rect()) == true {
				goto Send
			}
			if cs.clicktotype {
				goto Change
			}
			if !ms.Any() {
				goto Send
			}
		Change:
			/* change focus */
			if cs.focus != nil {
				cs.focus.Ctl("focus 0")
			}
			cs.focus = nil
			for _, c := range cs.active {
				if c.Hidden() == false && ms.Pt().In(c.Rect()) == true {
					cs.focus = c
					c.Ctl("focus 1")
					c.Mouse(ms)
					break
				}
			}
		Send:
			/* send mouse */
			if cs.focus != nil {
				cs.focus.Mouse(ms)
			}
			prevbut = ms.Buttons()
		case s := <-cs.ctl:
			cs.onctl(s)
		case <-cs.csexitc:
			/* fixme */
			return
		}
	}
}
