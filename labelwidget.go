package main

import (
	tb "github.com/nsf/termbox-go"
)

type LabelWidget struct {
	Rect     Rect
	Margin   Margin
	Color    Color
	Text     string
	Settings WidgetSetting
}

func NewLabelWidget(rect Rect, margin Margin, color Color, text string, settings WidgetSetting) *LabelWidget {
	if rect.H == 0 {
		rect.H = 1
	}
	if rect.W == 0 {
		rect.W = 10
	}

	InitColor(&color)

	w := LabelWidget{
		Rect:     rect,
		Margin:   margin,
		Color:    color,
		Text:     text,
		Settings: settings,
	}
	return &w
}

func (w *LabelWidget) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&FmtBox != 0 {
		boxRect := AddRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	rect := AddRectMargin(w.Rect, w.Margin)
	printRect(w.Text, rect, w.Color.Fg, w.Color.Bg)
}

func (w *LabelWidget) HandleEvent(e tb.Event) bool {
	return false
}

func (w *LabelWidget) SetText(text string) {
	w.Text = text
}
