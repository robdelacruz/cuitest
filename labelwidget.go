package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxLabel struct {
	Rect     TxRect
	Margin   TxMargin
	Color    TxColor
	Text     string
	Settings TxSetting
}

func NewTxLabel(rect TxRect, margin TxMargin, color TxColor, text string, settings TxSetting) *TxLabel {
	if rect.H == 0 {
		rect.H = 1
	}
	if rect.W == 0 {
		rect.W = 10
	}

	initColor(&color)

	w := TxLabel{
		Rect:     rect,
		Margin:   margin,
		Color:    color,
		Text:     text,
		Settings: settings,
	}
	return &w
}

func (w *TxLabel) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&TxFmtBox != 0 {
		boxRect := addRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	rect := addRectMargin(w.Rect, w.Margin)
	printRect(w.Text, rect, w.Color.Fg, w.Color.Bg)
}

func (w *TxLabel) HandleEvent(e tb.Event) bool {
	return false
}

func (w *TxLabel) SetText(text string) {
	w.Text = text
}
