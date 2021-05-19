package main

import (
	"fmt"
	"log"
	"os"

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
	labelw  *LabelWidget
	popupw  Widget
}

func NewMainWindow() *MainWindow {
	w := MainWindow{}
	w.width = _termW
	w.height = _termH

	w.smileyw = &Smiley{}
	w.smileyw.X = 5
	w.smileyw.Y = 5

	grey39 := tb.Attribute(242)
	gold1 := tb.Attribute(221)

	color := Color{
		Fg: gold1,
		Bg: grey39,
	}

	lbltext := "Now is the time for all good men to come to the aid of the party.\n\nThe quick brown fox jumps over the lazy dog. Now is the time for all good men to come to the aid of the party. The quick brown fox jumps over the lazy dog. Now is the time for all good men to come to the aid of the party. The quick brown fox jumps over the lazy dog."

	w.labelw = NewLabelWidget(Rect{1, 20, 30, 10}, MarginX, color, lbltext, FmtBox)

	return &w
}

func (w *MainWindow) popupCB(we *WidgetEvent) {
	if we.Code == WidgetEventEnter {
		_, ok := w.popupw.(*MenuWidget)
		if ok {
			sel := we.P1
			w.labelw.SetText(fmt.Sprintf("Selected menu option:\n%d", sel))
		}
		_, ok = w.popupw.(*ListboxWidget)
		if ok {
			selstr := we.Pstr
			w.labelw.SetText(fmt.Sprintf("Selected listbox item:\n%s", selstr))
		}
		w.popupw = nil
	} else if we.Code == WidgetEventEsc {
		w.labelw.SetText("*** Canceled operation ***")
		w.popupw = nil
	} else if we.Code == WidgetEventSel {
		w.labelw.SetText(we.Pstr)
	}
}

func (w *MainWindow) Draw() {
	tb.Clear(0, 0)

	w.smileyw.Draw()
	w.labelw.Draw()

	if w.popupw != nil {
		w.popupw.Draw()
	}

	tb.Flush()
}

func (w *MainWindow) HandleEvent(e tb.Event) bool {
	if w.popupw != nil {
		return w.popupw.HandleEvent(e)
	}

	white := tb.Attribute(16)
	black := tb.Attribute(17)
	darkolivegreen := tb.Attribute(156)
	darkorange := tb.Attribute(167)
	grey39 := tb.Attribute(242)
	plum1 := tb.Attribute(220)
	gold1 := tb.Attribute(221)

	color1 := Color{
		Fg:          darkolivegreen,
		Bg:          black,
		HighlightFg: grey39,
		HighlightBg: white,
	}
	color2 := Color{
		Fg:          darkorange,
		Bg:          black,
		HighlightFg: gold1,
		HighlightBg: grey39,
	}
	color3 := Color{
		Fg: plum1,
		Bg: grey39,
	}
	color4 := Color{
		Fg: darkolivegreen,
		Bg: grey39,
	}

	if e.Ch == 'm' {
		items := []*WidgetItem{
			&WidgetItem{1, "option1", "Menu Option 1 abc"},
			&WidgetItem{2, "option2", "Option 2 def"},
			&WidgetItem{3, "option3", "Option 3 ghijkl"},
			&WidgetItem{4, "option4", "Option 4 some more text"},
			&WidgetItem{5, "option5", "Option 5 xyz"},
			&WidgetItem{6, "option6", "Option 6 lmnop"},
			&WidgetItem{7, "option7", "Option 7 qrstuvw"},
			&WidgetItem{8, "option8", "Option 8 12345"},
			&WidgetItem{9, "option9", "Option 9 123"},
			&WidgetItem{10, "option10", "Option 10"},
		}
		//w.popupw = NewMenuWidget(Rect{5, 1, 0, 0}, Margin0, color1, w.popupCB, items, FmtBox|FmtCenter)
		w.popupw = NewMenuWidget(Rect{5, 1, 31, 0}, MarginX, color1, w.popupCB, items, FmtBox|FmtCenter)
		return true
	} else if e.Ch == 'l' {
		items := []*WidgetItem{
			&WidgetItem{1, "line1", "Now is the time"},
			&WidgetItem{2, "line2", "for all good men"},
			&WidgetItem{3, "line3", "to come to the aid"},
			&WidgetItem{4, "line4", "of the party."},
			&WidgetItem{5, "line5", "-- typing drill"},
		}
		w.popupw = NewListboxWidget(Rect{10, 1, 30, 0}, MarginX, color2, w.popupCB, items, FmtBox)
		return true
	} else if e.Ch == 't' {
		var cellColor Color
		cols := []CellSetting{
			CellSetting{0, 15, cellColor},
			CellSetting{15, 10, cellColor},
			CellSetting{35, 10, cellColor},
		}
		headings := []string{"col1", "col2", "col3"}
		rows := []TableRow{
			TableRow{"abc", "defghi", "jklmn"},
			TableRow{"ABC", "DEFGHI", "JKLMN"},
			TableRow{"12345", "678", "9012"},
			TableRow{"Now is", "the time", "for all"},
			TableRow{"good men", "to come to", "the aid"},
			TableRow{"of the", "party.", ""},
			TableRow{"12345", "678", "9012"},
		}
		//w.popupw = NewTableWidget(Rect{5, 5, 0, 0}, color3, color1, w.popupCB, cols, headings, rows, FmtBox)
		w.popupw = NewTableWidget(Rect{5, 5, 0, 0}, Margin1, color3, color1, w.popupCB, cols, headings, rows, FmtCenter)
		return true
	} else if e.Ch == 'e' {
		w.popupw = NewEntryWidget(Rect{5, 5, 30, 1}, Margin0, color4, w.popupCB, "Entry Text", FmtBox)
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
