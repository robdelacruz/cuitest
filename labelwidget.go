package main

import (
	tb "github.com/nsf/termbox-go"
)

type LabelWidgetSettings uint64

type LabelWidget struct {
	Rect     Rect
	Fg, Bg   tb.Attribute
	Text     string
	Settings LabelWidgetSettings
}

func NewLabelWidget(rect Rect, fg, bg tb.Attribute, text string, settings LabelWidgetSettings) *LabelWidget {
	if rect.H == 0 {
		rect.H = 1
	}
	if rect.W == 0 {
		// Add 1 char margin to the left and right of text.
		rect.W = len(text) + 2
	}

	// Truncate label text that go beyond width.
	if len(text)+2 > rect.W {
		text = text[:rect.W-2]
	}

	w := LabelWidget{
		Rect:     rect,
		Fg:       fg,
		Bg:       bg,
		Text:     text,
		Settings: settings,
	}
	return &w
}

func (w *LabelWidget) Draw() {
	clearRect(w.Rect, w.Bg)
	print(w.Text, w.Rect.X+1, w.Rect.Y, w.Fg, w.Bg)
}

func (w *LabelWidget) HandleEvent(e tb.Event) bool {
	return false
}

func (w *LabelWidget) SetText(text string) {
	w.Text = text
}
