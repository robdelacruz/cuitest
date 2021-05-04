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

func print(s string, x, y int, fg, bg tb.Attribute) {
	for _, c := range s {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
}

func print0(s string, x, y int) {
	for _, c := range s {
		tb.SetChar(x, y, c)
		x++
	}
}

func printcenter(s string, x, y int, fg, bg tb.Attribute, w int) {
	if w > len(s) {
		x += w/2 - len(s)/2
	}
	for _, c := range s {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
}

func clearRect(rect Rect, bg tb.Attribute) {
	for y := rect.y; y < rect.y+rect.h; y++ {
		for x := rect.x; x < rect.x+rect.w; x++ {
			tb.SetCell(x, y, ' ', 0, bg)
		}
	}
}

type Widget interface {
	Draw()
	HandleEvent(e tb.Event) bool
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
		"Option 1",
		"Option 2",
		"Option 3",
		"Option 4",
		"Option 5",
		"Option 6",
		"Option 7",
		"Option 8",
		"Option 9",
		"Option 10",
	}
	w.menu = NewMenuWidget(Rect{10, 10, 0, 4}, items, tb.Attribute(156), tb.Attribute(17))

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

type Pos struct {
	x, y int
}
type Rect struct {
	x, y, w, h int
}

type MenuWidget struct {
	rect    Rect
	items   []string
	isel    int
	fg, bg  tb.Attribute
	iscroll int
}

func NewMenuWidget(rect Rect, items []string, fg, bg tb.Attribute) *MenuWidget {
	// If not specified, automatically set width and height based on menu items.
	if rect.h == 0 {
		rect.h = len(items)
	}
	if rect.w == 0 {
		maxlen := 0
		for _, item := range items {
			if len(item) > maxlen {
				maxlen = len(item)
			}
		}
		// Add 1 char margin to the left and right of item.
		rect.w = maxlen + 2
	}

	// Remove menu items that go beyond height.
	//	if len(items) > rect.h {
	//		items = items[:rect.h]
	//	}

	// Truncate menu item text that go beyond width.
	for i, item := range items {
		if len(item) > rect.w {
			items[i] = item[:rect.w]
		}
	}

	// Pad menu items with spaces to fill width.
	//	for i, item := range items {
	//		if rect.w > len(item) {
	//			items[i] = item + strings.Repeat(" ", rect.w-len(item))
	//		}
	//	}

	w := MenuWidget{
		rect:    rect,
		items:   items,
		isel:    0,
		fg:      fg,
		bg:      bg,
		iscroll: 0,
	}
	return &w
}

func (w *MenuWidget) Draw() {
	clearRect(w.rect, w.bg)

	starti := w.iscroll
	endi := w.iscroll + w.rect.h - 1
	if endi > len(w.items)-1 {
		endi = len(w.items) - 1
	}

	y := w.rect.y
	for i := starti; i <= endi; i++ {
		item := w.items[i]
		if w.isel == i {
			// Highlight selected menu item
			printcenter(item, w.rect.x, y, w.bg, w.fg, w.rect.w)
		} else {
			printcenter(item, w.rect.x, y, w.fg, w.bg, w.rect.w)
		}
		y++
	}
}

func (w *MenuWidget) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch != 0 {
		return false
	}

	switch e.Key {
	case tb.KeyArrowUp:
		w.isel--
		if w.isel < 0 {
			w.isel = len(w.items) - 1
		}
		w.AdjustScroll()
		return true
	case tb.KeyArrowDown:
		w.isel++
		if w.isel > len(w.items)-1 {
			w.isel = 0
		}
		w.AdjustScroll()
		return true
	}
	return false
}

func (w *MenuWidget) AdjustScroll() {
	starti := w.iscroll
	endi := w.iscroll + w.rect.h - 1

	if w.isel < starti {
		w.iscroll -= starti - w.isel
	} else if w.isel > endi {
		w.iscroll += w.isel - endi
	}

	if w.iscroll < 0 {
		w.iscroll = 0
	} else if w.iscroll > len(w.items)-1 {
		w.iscroll = len(w.items) - 1
	}
}
