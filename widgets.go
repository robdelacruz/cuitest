package main

import (
	"strings"

	tb "github.com/nsf/termbox-go"
)

type Pos struct {
	x, y int
}
type Rect struct {
	x, y, w, h int
}

type Widget interface {
	Draw()
	HandleEvent(e tb.Event) bool
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

func drawBox(rect Rect, fg, bg tb.Attribute) {
	print("┌", rect.x, rect.y, fg, bg)
	print("┐", rect.x+rect.w-1, rect.y, fg, bg)

	hline := strings.Repeat("─", rect.w-2)
	print(hline, rect.x+1, rect.y, fg, bg)
	print(hline, rect.x+1, rect.y+rect.h-1, fg, bg)

	vchar := "│"
	for j := rect.y + 1; j < rect.y+rect.h-1; j++ {
		print(vchar, rect.x, j, fg, bg)
	}
	for j := rect.y + 1; j < rect.y+rect.h-1; j++ {
		print(vchar, rect.x+rect.w-1, j, fg, bg)
	}

	print("┘", rect.x+rect.w-1, rect.y+rect.h-1, fg, bg)
	print("└", rect.x, rect.y+rect.h-1, fg, bg)
}

type MenuWidget struct {
	rect    Rect
	items   []string
	isel    int
	fg, bg  tb.Attribute
	o       MenuWidgetSettings
	iscroll int
}

type MenuWidgetSettings uint64

const (
	MenuWidgetNormal MenuWidgetSettings = 1 << iota
	MenuWidgetCenter
	MenuWidgetBox
)

func NewMenuWidget(rect Rect, items []string, fg, bg tb.Attribute, o MenuWidgetSettings) *MenuWidget {
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

	// Truncate menu item text that go beyond width.
	for i, item := range items {
		if len(item)+2 > rect.w {
			items[i] = item[:rect.w-2]
		}
	}

	w := MenuWidget{
		rect:    rect,
		items:   items,
		isel:    0,
		fg:      fg,
		bg:      bg,
		o:       o,
		iscroll: 0,
	}
	return &w
}

func (w *MenuWidget) Draw() {
	clearRect(w.rect, w.bg)

	if w.o&MenuWidgetBox != 0 {
		boxRect := Rect{w.rect.x - 1, w.rect.y - 1, w.rect.w + 2, w.rect.h + 2}
		drawBox(boxRect, w.fg, w.bg)
	}

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
			if w.o&MenuWidgetCenter != 0 {
				printcenter(item, w.rect.x, y, w.bg, w.fg, w.rect.w)
			} else {
				print(item, w.rect.x+1, y, w.bg, w.fg)
			}
		} else {
			if w.o&MenuWidgetCenter != 0 {
				printcenter(item, w.rect.x, y, w.fg, w.bg, w.rect.w)
			} else {
				print(item, w.rect.x+1, y, w.fg, w.bg)
			}
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
