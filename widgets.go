package main

import (
	"strings"

	tb "github.com/nsf/termbox-go"
)

type Pos struct {
	X, Y int
}
type Rect struct {
	X, Y, W, H int
}

type Widget interface {
	Draw()
	HandleEvent(e tb.Event) bool
}

type WidgetEventCode int

const (
	WidgetEventEnter WidgetEventCode = iota
	WidgetEventEsc
	WidgetEventSel
)

type WidgetEvent struct {
	Code   WidgetEventCode
	P1     int
	P2     int
	Pnum   float32
	Pstr   string
	Detail interface{}
}
type WidgetEventCB func(we *WidgetEvent)

type WidgetSetting uint64

const (
	WidgetNormal WidgetSetting = 1 << iota
	WidgetCenter
	WidgetBox
)

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

func printspaces(nspace, x, y int, fg, bg tb.Attribute) {
	for i := 0; i < nspace; i++ {
		tb.SetCell(x+i, y, ' ', fg, bg)
	}
}

func printpadded(s string, nleftspace, nrightspace int, x, y int, fg, bg tb.Attribute) {
	print(strings.Repeat(" ", nleftspace), x, y, fg, bg)
	print(s, x+nleftspace, y, fg, bg)
	print(strings.Repeat(" ", nrightspace), x+nleftspace+len(s), y, fg, bg)
}

func printpaddedcenter(s string, x, y int, fg, bg tb.Attribute, w int) {
	xcenter := x
	if w > len(s) {
		xcenter += w/2 - len(s)/2
	}
	nleft := xcenter - x
	nright := w - nleft - len(s)
	printpadded(s, nleft, nright, x, y, fg, bg)
}

func clearRect(rect Rect, bg tb.Attribute) {
	for y := rect.Y; y < rect.Y+rect.H; y++ {
		for x := rect.X; x < rect.X+rect.W; x++ {
			tb.SetCell(x, y, ' ', 0, bg)
		}
	}
}

func drawBox(rect Rect, fg, bg tb.Attribute) {
	print("┌", rect.X, rect.Y, fg, bg)
	print("┐", rect.X+rect.W-1, rect.Y, fg, bg)

	hline := strings.Repeat("─", rect.W-2)
	print(hline, rect.X+1, rect.Y, fg, bg)
	print(hline, rect.X+1, rect.Y+rect.H-1, fg, bg)

	vchar := "│"
	for j := rect.Y + 1; j < rect.Y+rect.H-1; j++ {
		print(vchar, rect.X, j, fg, bg)
	}
	for j := rect.Y + 1; j < rect.Y+rect.H-1; j++ {
		print(vchar, rect.X+rect.W-1, j, fg, bg)
	}

	print("┘", rect.X+rect.W-1, rect.Y+rect.H-1, fg, bg)
	print("└", rect.X, rect.Y+rect.H-1, fg, bg)
}
