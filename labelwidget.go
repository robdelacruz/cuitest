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
		rect.W = 10
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
	//	print(w.Text, w.Rect.X+1, w.Rect.Y, w.Attrs.Fg, w.Attrs.Bg)
	printRect(w.Text, w.Rect, w.Attrs.Fg, w.Attrs.Bg)
}

func (w *LabelWidget) HandleEvent(e tb.Event) bool {
	return false
}

func (w *LabelWidget) SetText(text string) {
	w.Text = text
}
