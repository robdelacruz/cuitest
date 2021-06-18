package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxLabelEntry struct {
	Rect       TxRect
	Margin     TxMargin
	LabelColor TxColor
	Color      TxColor
	Cb         TxEventCB
	LabelText  string
	Text       string
	Settings   TxSetting
	Cur        TxPos
}

func NewTxLabelEntry(rect TxRect, margin TxMargin, labelcolor, color TxColor, cb TxEventCB, labeltext, text string, settings TxSetting) *TxLabelEntry {
	if rect.H < 2 {
		rect.H = 2
	}
	if rect.W == 0 {
		rect.W = 10
	}

	initColor(&labelcolor)
	initColor(&color)

	w := TxLabelEntry{
		Rect:       rect,
		Margin:     margin,
		LabelColor: labelcolor,
		Color:      color,
		Cb:         cb,
		LabelText:  labeltext,
		Text:       text,
		Settings:   settings,
		Cur:        TxPos{0, 0},
	}
	return &w
}

func (w *TxLabelEntry) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&TxFmtBox != 0 {
		boxRect := addRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	rect := addRectMargin(w.Rect, w.Margin)
	y := rect.Y
	printw(w.LabelText, rect.X, y, w.LabelColor.Fg, w.LabelColor.Bg, rect.W)
	y++
	printw(w.Text, rect.X, y, w.Color.Fg, w.Color.Bg, rect.W)

	// Print cursor
	curCh := ' '
	if w.Cur.X <= len(w.Text)-1 {
		curCh = rune(w.Text[w.Cur.X])
	}
	tb.SetCell(rect.X+w.Cur.X, y, curCh, w.Color.HighlightFg, w.Color.HighlightBg)
}

func (w *TxLabelEntry) HandleEvent(e tb.Event) bool {
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

func (w *TxLabelEntry) adjustCur() {
	if w.Cur.X < 0 {
		w.Cur.X = 0
	} else if w.Cur.X > len(w.Text) {
		w.Cur.X = len(w.Text)
	}
}

func (w *TxLabelEntry) InsertChar(r rune, x int) {
	if x > len(w.Text)-1 {
		w.Text += string(r)
	} else {
		w.Text = w.Text[:x] + string(r) + w.Text[x:]
	}
}

func (w *TxLabelEntry) SetText(text string) {
	w.Text = text
}
