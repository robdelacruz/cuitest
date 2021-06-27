package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxLabel struct {
	Props *TxProps
	Text  string
}

func NewTxLabel(props *TxProps, text string) *TxLabel {
	if props == nil {
		props = defaultProps()
	}

	if props.Rect.H == 0 {
		props.Rect.H = 1
	}
	if props.Rect.W == 0 {
		props.Rect.W = 10
	}

	initColor(&props.Clr)

	w := TxLabel{
		Props: props,
		Text:  text,
	}
	return &w
}

func (w *TxLabel) Draw() {
	p := w.Props
	clearRect(p.Rect, p.Clr.Bg)

	if p.Fmt&TxFmtBox != 0 {
		boxRect := addRectBox(p.Rect)
		drawBox(boxRect, p.Clr.Fg, p.Clr.Bg)
	}

	rect := addRectMargin(p.Rect, p.Margin)
	printRect(w.Text, rect, p.Clr.Fg, p.Clr.Bg)
}

func (w *TxLabel) HandleEvent(e tb.Event) bool {
	return false
}

func (w *TxLabel) SetText(text string) {
	w.Text = text
}
