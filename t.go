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

	smileyw *Smiley
	popupw  Widget
}

func NewMainWindow() *MainWindow {
	w := MainWindow{}
	w.width = _termW
	w.height = _termH

	w.smileyw = &Smiley{}
	w.smileyw.X = 5
	w.smileyw.Y = 5

	return &w
}

func (w *MainWindow) popupCB(we *WidgetEvent) {
	if we.Code == WidgetEventEnter {
		w.popupw = nil
	} else if we.Code == WidgetEventEsc {
		w.popupw = nil
	}
}

func (w *MainWindow) Draw() {
	tb.Clear(0, 0)

	w.smileyw.Draw()

	if w.popupw != nil {
		w.popupw.Draw()
	}

	tb.Flush()
}

func (w *MainWindow) HandleEvent(e tb.Event) bool {
	if w.popupw != nil {
		return w.popupw.HandleEvent(e)
	}

	if e.Ch == 'm' {
		items := []string{
			"Menu Option 1 abc",
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
		w.popupw = NewMenuWidget(Rect{10, 10, 0, 0}, tb.Attribute(156), tb.Attribute(17), w.popupCB, items, MenuWidgetBox|MenuWidgetCenter)
		return true
	} else if e.Ch == 'l' {
		items := []string{
			"Now is the time",
			"for all good men",
			"to come to the aid",
			"of the party.",
			"-- typing drill",
		}
		w.popupw = NewListboxWidget(Rect{10, 20, 0, 0}, tb.Attribute(156), tb.Attribute(17), w.popupCB, items, ListboxWidgetBox)
		return true
	}

	switch e.Key {
	}

	return w.smileyw.HandleEvent(e)
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
