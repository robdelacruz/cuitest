package main

import (
	"regexp"
	"strings"

	tb "github.com/nsf/termbox-go"
)

type TxEntry struct {
	Rect         TxRect
	Margin       TxMargin
	Color        TxColor
	Cb           TxEventCB
	Text         string
	ValidatorReg *regexp.Regexp
	Settings     TxSetting
	Cur          TxPos
}

func NewTxEntry(rect TxRect, margin TxMargin, color TxColor, cb TxEventCB, text string, svalidator string, settings TxSetting) *TxEntry {
	if rect.H == 0 {
		rect.H = 1
	}
	if rect.W == 0 {
		rect.W = 10
	}

	initColor(&color)

	var validatorReg *regexp.Regexp
	if svalidator != "" {
		if !strings.HasPrefix(svalidator, "^") {
			svalidator = "^" + svalidator
		}
		if !strings.HasSuffix(svalidator, "$") {
			svalidator = svalidator + "$"
		}

		var err error
		validatorReg, err = regexp.Compile(svalidator)
		if err != nil {
			validatorReg = nil
		}
	}

	w := TxEntry{
		Rect:         rect,
		Margin:       margin,
		Color:        color,
		Cb:           cb,
		Text:         "",
		ValidatorReg: validatorReg,
		Settings:     settings,
		Cur:          TxPos{0, 0},
	}
	if text != "" {
		w.SetText(text)
	}
	return &w
}

func (w *TxEntry) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&TxFmtBox != 0 {
		boxRect := addRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	rect := addRectMargin(w.Rect, w.Margin)
	printw(w.Text, rect.X, rect.Y, w.Color.Fg, w.Color.Bg, rect.W)

	// Print cursor
	curCh := ' '
	if w.Cur.X <= len(w.Text)-1 {
		curCh = rune(w.Text[w.Cur.X])
	}
	tb.SetCell(rect.X+w.Cur.X, rect.Y, curCh, w.Color.HighlightFg, w.Color.HighlightBg)
}

func (w *TxEntry) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch != 0 {
		if w.InsertChar(e.Ch, w.Cur.X) {
			w.Cur.X++
		}
		return true
	}

	switch e.Key {
	case tb.KeySpace:
		if w.InsertChar(' ', w.Cur.X) {
			w.Cur.X++
		}
		return true
	case tb.KeyEsc:
		if w.Cb != nil {
			we := TxEvent{
				Code: TxEventEsc,
			}
			w.Cb(&we)
		}
		return true
	case tb.KeyArrowLeft:
		w.Cur.X--
		w.adjustCur()
		return true
	case tb.KeyArrowRight:
		w.Cur.X++
		w.adjustCur()
		return true
	case tb.KeyCtrlA:
		w.Cur.X = 0
		w.adjustCur()
		return true
	case tb.KeyCtrlE:
		w.Cur.X = len(w.Text)
		w.adjustCur()
		return true
	case tb.KeyBackspace:
		fallthrough
	case tb.KeyBackspace2:
		if w.Cur.X > 0 {
			w.Cur.X--
			w.Text = w.Text[:w.Cur.X] + w.Text[w.Cur.X+1:]
		}
		return true
	case tb.KeyDelete:
		if w.Cur.X <= len(w.Text)-1 {
			w.Text = w.Text[:w.Cur.X] + w.Text[w.Cur.X+1:]
		}
		return true
	}
	return false
}

func (w *TxEntry) adjustCur() {
	if w.Cur.X < 0 {
		w.Cur.X = 0
	} else if w.Cur.X > len(w.Text) {
		w.Cur.X = len(w.Text)
	}
}

func (w *TxEntry) InsertChar(r rune, x int) bool {
	var newText string
	if x > len(w.Text)-1 {
		newText = w.Text + string(r)
	} else {
		newText = w.Text[:x] + string(r) + w.Text[x:]
	}

	return w.SetText(newText)
}

func (w *TxEntry) SetText(text string) bool {
	if w.ValidatorReg == nil || w.ValidatorReg.MatchString(text) {
		w.Text = text
		return true
	}
	return false
}
