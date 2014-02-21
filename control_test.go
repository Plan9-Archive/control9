package control

import (
	"bitbucket.org/mischief/draw9"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

func init() {
	logf, err := os.Create("/tmp/control.log")
	if err != nil {
		panic(err)
	}

	log.SetOutput(logf)
}

var (
	disp *draw9.Display
	kbd  *draw9.Consctl
	ms   *draw9.Mousectl

	setup sync.Once
)

func initdraw() {
	d, err := draw9.InitDraw(nil, "", "controltest")
	if err != nil {
		panic(err)
	}
	disp = d

	kbd = draw9.InitCons("")
	ms = draw9.InitMouse("", d.ScreenImage)

	err = InitControlset(d)
	if err != nil {
		panic(err)
	}
}

func setupTest() {
	setup.Do(initdraw)
}

func TestButton(t *testing.T) {
	setupTest()

	cs, err := NewControlset(disp.ScreenImage, kbd, ms)

	if err != nil {
		t.Fatal(err)
	}

	row := cs.Row("row")
	e := cs.TextButton3("exit")
	cs.Ctl("exit text exit")
	cs.Ctl("row add exit")

	d := cs.TextButton3("doit")
	cs.Ctl("doit text doit")
	cs.Ctl("row add doit")

	cs.OnResize = func(cs *Controlset) {
		disp.Attach(draw9.Refnone)

		*cs.screen = *disp.ScreenImage
		cs.Ctl("row size")
		r := disp.ScreenImage.R
		cs.Ctl(fmt.Sprintf("row rect %d %d %d %d", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y))
		log.Printf("row %+v", row)
		cs.Ctl("row show")
		disp.Flush()
	}

	//cs.Ctl("row show")
	cs.Resize()

	evt := make(chan string, 1)
	e.Wire(evt)
	d.Wire(evt)

	for ev := range evt {

		t.Logf("evt: %s", ev)
		f := strings.Fields(ev)
		if strings.Contains(f[2], "4") {
			break
		}
	}
}
