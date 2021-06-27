package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxLabelEntry struct {
	Props *TxProps
	label *TxLabel
	entry *TxEntry
}

func NewTxLabelEntry(props *TxProps, labelclr, entryclr TxColor, labeltext, entrytext string, svalidator string) *TxLabelEntry {
	if props == nil {
		props = defaultProps()
	}

	if props.Rect.H < 2 {
		props.Rect.H = 2
	}
	if props.Rect.W == 0 {
		props.Rect.W = 10
	}

	initColor(&labelclr)
	initColor(&entryclr)

	rectmargin := addRectMargin(props.Rect, props.Margin)
	labelrect := TxRect{rectmargin.X, rectmargin.Y, rectmargin.W, rectmargin.H}
	labelprops := &TxProps{labelrect, TxMargin0, labelclr, nil, 0}
	label := NewTxLabel(labelprops, labeltext)

	entryrect := TxRect{rectmargin.X, rectmargin.Y + 1, rectmargin.W, rectmargin.H}
	entryprops := &TxProps{entryrect, TxMargin0, entryclr, props.EventCB, props.Fmt}
	entry := NewTxEntry(entryprops, entrytext, svalidator)

	w := TxLabelEntry{
		Props: props,
		label: label,
		entry: entry,
	}
	return &w
}

func (w *TxLabelEntry) Draw() {
	p := w.Props
	clearRect(p.Rect, p.Clr.Bg)

	if p.Fmt&TxFmtBox != 0 {
		boxRect := addRectBox(p.Rect)
		drawBox(boxRect, p.Clr.Fg, p.Clr.Bg)
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
