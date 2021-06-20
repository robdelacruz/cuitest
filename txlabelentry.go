package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxLabelEntry struct {
	Rect     TxRect
	Margin   TxMargin
	Color    TxColor
	Settings TxSetting
	label    *TxLabel
	entry    *TxEntry
}

func NewTxLabelEntry(rect TxRect, margin TxMargin, labelcolor, color TxColor, cb TxEventCB, labeltext, text string, svalidator string, settings TxSetting) *TxLabelEntry {
	if rect.H < 2 {
		rect.H = 2
	}
	if rect.W == 0 {
		rect.W = 10
	}

	initColor(&labelcolor)
	initColor(&color)

	rectmargin := addRectMargin(rect, margin)
	rectlabel := TxRect{rectmargin.X, rectmargin.Y, rectmargin.W, rectmargin.H}
	rectentry := TxRect{rectmargin.X, rectmargin.Y + 1, rectmargin.W, rectmargin.H}
	label := NewTxLabel(rectlabel, TxMargin0, labelcolor, labeltext, settings)
	entry := NewTxEntry(rectentry, TxMargin0, color, cb, text, svalidator, settings)

	w := TxLabelEntry{
		Rect:     rect,
		Margin:   margin,
		Color:    color,
		Settings: settings,
		label:    label,
		entry:    entry,
	}
	return &w
}

func (w *TxLabelEntry) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&TxFmtBox != 0 {
		boxRect := addRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	w.label.Draw()
	w.entry.Draw()
}

func (w *TxLabelEntry) HandleEvent(e tb.Event) bool {
	return w.entry.HandleEvent(e)
}

func (w *TxLabelEntry) SetText(text string) bool {
	return w.entry.SetText(text)
}
