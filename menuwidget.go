package main

import (
	tb "github.com/nsf/termbox-go"
)

type MenuWidget struct {
	Rect      Rect
	Margin    Margin
	Attrs     WidgetAttributes
	Cb        WidgetEventCB
	Items     []*WidgetItem
	Settings  WidgetSetting
	Sel       int
	Scrollpos int
}

func NewMenuWidget(rect Rect, margin Margin, attrs WidgetAttributes, cb WidgetEventCB, items []*WidgetItem, settings WidgetSetting) *MenuWidget {
	// If not specified, automatically set width and height based on menu items.
	if rect.H == 0 {
		rect.H = len(items) + margin.T + margin.B
	}
	if rect.W == 0 {
		maxlen := 0
		for _, item := range items {
			if len(item.Display) > maxlen {
				maxlen = len(item.Display)
			}
		}
		rect.W = maxlen + margin.L + margin.R
	}

	// Truncate menu item text that go beyond width.
	for i, item := range items {
		if len(item.Display) > rect.W-margin.L-margin.R {
			items[i].Display = item.Display[:rect.W-margin.L-margin.R]
		}
	}

	InitWidgetAttributes(&attrs)

	w := MenuWidget{
		Rect:      rect,
		Margin:    margin,
		Attrs:     attrs,
		Cb:        cb,
		Items:     items,
		Settings:  settings,
		Sel:       0,
		Scrollpos: 0,
	}
	return &w
}

func (w *MenuWidget) Draw() {
	clearRect(w.Rect, w.Attrs.Bg)

	if w.Settings&WidgetBox != 0 {
		boxRect := AddRectBox(w.Rect)
		drawBox(boxRect, w.Attrs.Fg, w.Attrs.Bg)
	}

	contentRect := AddRectMargin(w.Rect, w.Margin)

	starti := w.Scrollpos
	endi := w.Scrollpos + contentRect.H - 1
	if endi > len(w.Items)-1 {
		endi = len(w.Items) - 1
	}

	y := contentRect.Y
	for i := starti; i <= endi; i++ {
		item := w.Items[i]
		if w.Sel == i {
			// Highlight selected menu item
			printspaces(w.Rect.W, w.Rect.X, y, w.Attrs.HighlightFg, w.Attrs.HighlightBg)
			if w.Settings&WidgetCenter != 0 {
				printcenter(item.Display, contentRect.X, y, w.Attrs.HighlightFg, w.Attrs.HighlightBg, contentRect.W)
				printcenter(item.Display, w.Rect.X, y, w.Attrs.HighlightFg, w.Attrs.HighlightBg, w.Rect.W)
			} else {
				printw(item.Display, contentRect.X, y, w.Attrs.HighlightFg, w.Attrs.HighlightBg, contentRect.W)
			}
		} else {
			if w.Settings&WidgetCenter != 0 {
				printcenter(item.Display, contentRect.X, y, w.Attrs.Fg, w.Attrs.Bg, contentRect.W)
			} else {
				printw(item.Display, contentRect.X, y, w.Attrs.Fg, w.Attrs.Bg, contentRect.W)
			}
		}
		y++
	}
}

func (w *MenuWidget) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch != 0 {
		return false
	}

	switch e.Key {
	case tb.KeyArrowUp:
		w.Sel--
		if w.Sel < 0 {
			w.Sel = len(w.Items) - 1
		}
		w.AdjustScroll()
		return true
	case tb.KeyArrowDown:
		w.Sel++
		if w.Sel > len(w.Items)-1 {
			w.Sel = 0
		}
		w.AdjustScroll()
		return true
	case tb.KeyEnter:
		if w.Cb != nil {
			we := WidgetEvent{
				Code: WidgetEventEnter,
				P1:   w.Sel,
			}
			w.Cb(&we)
		}
		return true
	case tb.KeyEsc:
		if w.Cb != nil {
			we := WidgetEvent{
				Code: WidgetEventEsc,
			}
			w.Cb(&we)
		}
		return true
	}
	return false
}

func (w *MenuWidget) AdjustScroll() {
	rect := AddRectMargin(w.Rect, w.Margin)
	starti := w.Scrollpos
	endi := w.Scrollpos + rect.H - 1

	if w.Sel < starti {
		w.Scrollpos -= starti - w.Sel
	} else if w.Sel > endi {
		w.Scrollpos += w.Sel - endi
	}

	if w.Scrollpos < 0 {
		w.Scrollpos = 0
	} else if w.Scrollpos > len(w.Items)-1 {
		w.Scrollpos = len(w.Items) - 1
	}
}
