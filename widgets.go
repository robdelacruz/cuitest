package main

import (
	"strings"
	"unicode"

	tb "github.com/nsf/termbox-go"
)

type Pos struct {
	X, Y int
}
type Rect struct {
	X, Y, W, H int
}
type Margin struct {
	T, B, L, R int
}

var Margin0 = Margin{0, 0, 0, 0}
var Margin1 = Margin{1, 1, 1, 1}
var MarginX = Margin{0, 0, 1, 1}
var MarginY = Margin{1, 1, 0, 0}

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

type WidgetAttributes struct {
	Fg, Bg                   tb.Attribute
	HighlightFg, HighlightBg tb.Attribute
}

type WidgetItem struct {
	Id      int
	Sid     string
	Display string
}

func InitWidgetAttributes(attrs *WidgetAttributes) {
	if attrs.Fg == 0 {
		attrs.Fg = tb.ColorWhite
	}
	if attrs.Bg == 0 {
		attrs.Bg = tb.ColorBlack
	}
	if attrs.HighlightFg == 0 {
		attrs.HighlightFg = attrs.Bg
	}
	if attrs.HighlightBg == 0 {
		attrs.HighlightBg = attrs.Fg
	}
}

func AddRectMargin(rect Rect, m Margin) Rect {
	return Rect{rect.X + m.L, rect.Y + m.T, rect.W - m.L - m.R, rect.H - m.T - m.B}
}
func AddRectBox(rect Rect) Rect {
	return Rect{rect.X - 1, rect.Y - 1, rect.W + 2, rect.H + 2}
}

func print0(s string, x, y int) {
	for _, c := range s {
		tb.SetChar(x, y, c)
		x++
	}
}

func print(s string, x, y int, fg, bg tb.Attribute) {
	for _, c := range s {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
}

func printw(s string, x, y int, fg, bg tb.Attribute, w int) {
	for i, c := range s {
		if i > w-1 {
			return
		}
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

func printcenter(s string, x, y int, fg, bg tb.Attribute, w int) {
	if w > len(s) {
		x += w/2 - len(s)/2
	}
	for _, c := range s {
		tb.SetCell(x, y, c, fg, bg)
		x++
	}
}

func printRect(s string, rect Rect, fg, bg tb.Attribute) {
	ss := lexwords(s)

	y := rect.Y
	x := rect.X
	for _, word := range ss {
		if word == "\n" {
			y++
			x = rect.X
			continue
		}

		// If word exceeds rect, go to next line.
		if x+len(word)-1 > rect.X+rect.W-1 {
			y++
			x = rect.X
		}

		// Stop if no more rect space to write.
		if y > rect.Y+rect.H-1 {
			break
		}

		// Skip over any whitespace at the start of the line (except for the first row).
		if len(word) == 1 && unicode.IsSpace([]rune(word)[0]) && x == rect.X && y > rect.Y {
			continue
		}

		print(word, x, y, fg, bg)
		x += len(word)
	}
}

// "Here's a sentence." ==> ["Here's", " ", "a", " ", "sentence."]
func lexwords(s string) []string {
	ss := []string{}
	var sb strings.Builder

	for _, r := range s {
		if unicode.IsSpace(r) {
			if sb.Len() > 0 {
				ss = append(ss, sb.String())
				sb.Reset()
			}
			ss = append(ss, string(r))
			continue
		}

		sb.WriteRune(r)
	}

	if sb.Len() > 0 {
		ss = append(ss, sb.String())
	}
	return ss
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
