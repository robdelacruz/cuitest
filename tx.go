package main

import (
	"strings"
	"unicode"

	tb "github.com/nsf/termbox-go"
)

type TxPos struct {
	X, Y int
}
type TxRect struct {
	X, Y, W, H int
}
type TxMargin struct {
	T, B, L, R int
}

var TxMargin0 = TxMargin{0, 0, 0, 0}
var TxMargin1 = TxMargin{1, 1, 1, 1}
var TxMarginX = TxMargin{0, 0, 1, 1}
var TxMarginY = TxMargin{1, 1, 0, 0}

type TxWidget interface {
	Draw()
	HandleEvent(e tb.Event) bool
}

type TxEventCode int

const (
	TxEventEnter TxEventCode = iota
	TxEventEsc
	TxEventSel
)

type TxEvent struct {
	Code   TxEventCode
	Item   *TxItem
	Detail interface{}
}
type TxEventCB func(we *TxEvent)

type TxSetting uint64

const (
	TxFmtNormal TxSetting = 1 << iota
	TxFmtCenter
	TxFmtBox
)

type TxColor struct {
	Fg, Bg                   tb.Attribute
	HighlightFg, HighlightBg tb.Attribute
}

var TxWhite = tb.Attribute(16)
var TxBlack = tb.Attribute(17)
var TxDarkolivegreen = tb.Attribute(156)
var TxDarkorange = tb.Attribute(167)
var TxGrey39 = tb.Attribute(242)
var TxPlum1 = tb.Attribute(220)
var TxGold1 = tb.Attribute(221)

var TxColorBW = TxColor{tb.ColorWhite, tb.ColorBlack, tb.ColorBlack, tb.ColorWhite}
var TxColorWhite = TxColor{TxWhite, TxBlack, TxBlack, TxWhite}
var TxColorGreen = TxColor{TxDarkolivegreen, TxBlack, TxBlack, TxDarkolivegreen}

type TxItem struct {
	Id      int64
	Alias   string
	Display string
}

func initColor(color *TxColor) {
	if color.Fg == 0 {
		color.Fg = tb.ColorWhite
	}
	if color.Bg == 0 {
		color.Bg = tb.ColorBlack
	}
	if color.HighlightFg == 0 {
		color.HighlightFg = color.Bg
	}
	if color.HighlightBg == 0 {
		color.HighlightBg = color.Fg
	}
}

func addRectMargin(rect TxRect, m TxMargin) TxRect {
	return TxRect{rect.X + m.L, rect.Y + m.T, rect.W - m.L - m.R, rect.H - m.T - m.B}
}
func addRectBox(rect TxRect) TxRect {
	return TxRect{rect.X - 1, rect.Y - 1, rect.W + 2, rect.H + 2}
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

func printRect(s string, rect TxRect, fg, bg tb.Attribute) {
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

func clearRect(rect TxRect, bg tb.Attribute) {
	for y := rect.Y; y < rect.Y+rect.H; y++ {
		for x := rect.X; x < rect.X+rect.W; x++ {
			tb.SetCell(x, y, ' ', 0, bg)
		}
	}
}

func drawBox(rect TxRect, fg, bg tb.Attribute) {
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
