package main

import (
	"log"
	"os"
	//	"strings"

	tb "github.com/nsf/termbox-go"
)

var _log *log.Logger
var _termW, _termH int

func main() {
	flog, err := os.Create("./log.txt")
	if err != nil {
		panic(err)
	}
	defer flog.Close()
	_log = log.New(flog, "", 0)

	err = tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()
	tb.SetOutputMode(tb.Output256)

	_termW, _termH = tb.Size()
	mainwin := NewMainWindow()
	mainwin.Draw()

	chev := make(chan tb.Event)
	go pollEvents(chev)

	for {
		e := <-chev

		if e.Ch == 'q' {
			break
		}
		if mainwin.HandleEvent(e) {
			mainwin.Draw()
		}
	}
}

func pollEvents(chev chan tb.Event) {
	for {
		e := tb.PollEvent()
		if e.Type != tb.EventKey {
			continue
		}
		chev <- e
	}
}

type MainWindow struct {
	width, height int
	isMenuActive  bool

	smiley       *Smiley
	menu         *MenuWidget
	activeWidget Widget
}

func NewMainWindow() *MainWindow {
	w := MainWindow{}
	w.width = _termW
	w.height = _termH

	w.smiley = &Smiley{}
	w.smiley.X = 5
	w.smiley.Y = 5

	items := []string{
		"Option 1 abc",
		"Option 2 def",
		"Option 3 ghijkl",
		"Option 4 some more text",
		"Option 5 xyz",
		"Option 6 lmnop",
		"Option 7 qrstuvw",
		"Option 8 12345",
		"Option 9 123",
		"Option 10",
	}
	//w.menu = NewMenuWidget(Rect{10, 10, 0, 4}, items, tb.Attribute(156), tb.Attribute(17), MenuWidgetCenter)
	//w.menu = NewMenuWidget(Rect{10, 10, 0, 4}, items, tb.Attribute(156), tb.Attribute(17), 0)
	w.menu = NewMenuWidget(Rect{10, 10, 0, 7}, items, tb.Attribute(156), tb.Attribute(17), MenuWidgetBox|MenuWidgetCenter)
	//w.menu = NewMenuWidget(Rect{10, 10, 0, 4}, items, tb.Attribute(156), tb.Attribute(17), MenuWidgetNormal)

	w.isMenuActive = false
	w.activeWidget = w.smiley

	return &w
}

func (w *MainWindow) Draw() {
	tb.Clear(0, 0)

	w.smiley.Draw()

	if w.isMenuActive {
		w.menu.Draw()
	}

	tb.Flush()
}

func (w *MainWindow) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}

	switch e.Key {
	case tb.KeyCtrlM:
		w.isMenuActive = !w.isMenuActive
		if w.isMenuActive {
			w.activeWidget = w.menu
		} else {
			w.activeWidget = w.smiley
		}
		return true
	}

	if w.activeWidget != nil {
		return w.activeWidget.HandleEvent(e)
	}
	return true
}

type Smiley struct {
	X, Y int
}

func (sm *Smiley) Draw() {
	fg := tb.Attribute(156)
	bg := tb.Attribute(0)
	tb.SetCell(sm.X, sm.Y, '😀', fg, bg)
}

func (sm *Smiley) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch != 0 {
		return false
	}

	switch e.Key {
	case tb.KeyArrowUp:
		if sm.Y > 0 {
			sm.Y--
			return true
		}
	case tb.KeyArrowDown:
		if sm.Y < _termH-1 {
			sm.Y++
			return true
		}
	case tb.KeyArrowLeft:
		if sm.X > 0 {
			sm.X--
			return true
		}
	case tb.KeyArrowRight:
		if sm.X < _termW-1 {
			sm.X++
			return true
		}
	}
	return false
}
