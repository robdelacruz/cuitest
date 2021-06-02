package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type TxCell struct {
	Rect     TxRect
	clr      TxColor
	sfmt     string
	v        interface{}
	settings TxSetting
}

func NewTxCell(rect TxRect, clr TxColor, sfmt string, v interface{}, settings TxSetting) *TxCell {
	if rect.H == 0 {
		rect.H = 1
	}
	if rect.W == 0 {
		rect.W = 10
	}

	initColor(&clr)

	w := TxCell{
		Rect:     rect,
		clr:      clr,
		sfmt:     sfmt,
		v:        v,
		settings: settings,
	}
	return &w
}

func (w *TxCell) Draw() {
	clearRect(w.Rect, w.clr.Bg)

	if w.settings&TxFmtBox != 0 {
		boxRect := addRectBox(w.Rect)
		drawBox(boxRect, w.clr.Fg, w.clr.Bg)
	}

	s := fmt.Sprintf(w.sfmt, w.v)
	if w.settings&TxFmtCenter != 0 {
		printcenter(s, w.Rect.X, w.Rect.Y, w.clr.Fg, w.clr.Bg, w.Rect.W)
	} else {
		printw(s, w.Rect.X, w.Rect.Y, w.clr.Fg, w.clr.Bg, w.Rect.W)
	}
}

func (w *TxCell) HandleEvent(e tb.Event) bool {
	return false
}

func (w *TxCell) SetFmt(sfmt string) {
	w.sfmt = sfmt
}
func (w *TxCell) SetVal(v interface{}) {
	w.v = v
}
