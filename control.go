package control

import (
	"bitbucket.org/mischief/draw9"
	"image"
)

type ControlType int

const (
	CtlUnknown ControlType = iota
	CtlBox
	CtlButton
	CtlEntry
	CtlRow
	CtlTextButton3
)

type Control interface {
	// Per-control functions
	Ctl(e string)
	Mouse(m draw9.Mouse)
	Key(c rune)
	Exit()
	SetSize()
	Activate()

	// BaseControl functions
	Name() string
	Rect() image.Rectangle
	Size() image.Rectangle
	Wire(e chan string)
	Hidden() bool
}

type BaseControl struct {
	name       string
	rect       image.Rectangle
	size       image.Rectangle
	event      chan string
	data       chan string
	typ        ControlType
	hidden     bool
	set        *Controlset
	screen     *draw9.Image
	format     string
	nextactive *Control
	next       *Control
}

func (bc *BaseControl) Name() string {
	return bc.name
}

func (bc *BaseControl) Rect() image.Rectangle {
	return bc.rect
}

func (bc *BaseControl) Size() image.Rectangle {
	return bc.size
}

func (bc *BaseControl) Wire(e chan string) {
	bc.event = e
}

func (bc *BaseControl) Hidden() bool {
	return bc.hidden
}
