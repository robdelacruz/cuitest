package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxEntry struct {
	Rect     TxRect
	Margin   TxMargin
	Color    TxColor
	Cb       TxEventCB
	Text     string
	Settings TxSetting
	Cur      TxPos
}

func NewTxEntry(rect TxRect, margin TxMargin, color TxColor, cb TxEventCB, text string, settings TxSetting) *TxEntry {
	if rect.H == 0 {
		rect.H = 1
	}
	if rect.W == 0 {
		rect.W = 10
	}

	initColor(&color)

	w := TxEntry{
		Rect:     rect,
		Margin:   margin,
		Color:    color,
		Cb:       cb,
		Text:     text,
		Settings: settings,
		Cur:      TxPos{0, 0},
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
		w.InsertChar(e.Ch, w.Cur.X)
		w.Cur.X++
		return true
	}

	switch e.Key {
	case tb.KeySpace:
		w.InsertChar(' ', w.Cur.X)
		w.Cur.X++
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

func (w *TxEntry) InsertChar(r rune, x int) {
	if x > len(w.Text)-1 {
		w.Text += string(r)
	} else {
		w.Text = w.Text[:x] + string(r) + w.Text[x:]
	}
}

func (w *TxEntry) SetText(text string) {
	w.Text = text
}
