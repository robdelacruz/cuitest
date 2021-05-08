package main

import (
	tb "github.com/nsf/termbox-go"
)

type LabelWidget struct {
	Rect     Rect
	Attrs    WidgetAttributes
	Text     string
	Settings WidgetSetting
}

func NewLabelWidget(rect Rect, attrs WidgetAttributes, text string, settings WidgetSetting) *LabelWidget {
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

	InitWidgetAttributes(&attrs)

	w := LabelWidget{
		Rect:     rect,
		Attrs:    attrs,
		Text:     text,
		Settings: settings,
	}
	return &w
}

func (w *LabelWidget) Draw() {
	clearRect(w.Rect, w.Attrs.Bg)
	print(w.Text, w.Rect.X+1, w.Rect.Y, w.Attrs.Fg, w.Attrs.Bg)
}

func (w *LabelWidget) HandleEvent(e tb.Event) bool {
	return false
}

func (w *LabelWidget) SetText(text string) {
	w.Text = text
}
