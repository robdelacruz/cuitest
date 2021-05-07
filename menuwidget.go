package main

import (
	tb "github.com/nsf/termbox-go"
)

type MenuWidget struct {
	Rect      Rect
	Fg, Bg    tb.Attribute
	Cb        WidgetEventCB
	Items     []string
	Settings  WidgetSetting
	Sel       int
	Scrollpos int
}

func NewMenuWidget(rect Rect, fg, bg tb.Attribute, cb WidgetEventCB, items []string, settings WidgetSetting) *MenuWidget {
	// If not specified, automatically set width and height based on menu items.
	if rect.H == 0 {
		rect.H = len(items)
	}
	if rect.W == 0 {
		maxlen := 0
		for _, item := range items {
			if len(item) > maxlen {
				maxlen = len(item)
			}
		}
		// Add 1 char margin to the left and right of item.
		rect.W = maxlen + 2
	}

	// Truncate menu item text that go beyond width.
	for i, item := range items {
		if len(item)+2 > rect.W {
			items[i] = item[:rect.W-2]
		}
	}

	w := MenuWidget{
		Rect:      rect,
		Fg:        fg,
		Bg:        bg,
		Cb:        cb,
		Items:     items,
		Settings:  settings,
		Sel:       0,
		Scrollpos: 0,
	}
	return &w
}

func (w *MenuWidget) Draw() {
	clearRect(w.Rect, w.Bg)

	if w.Settings&WidgetBox != 0 {
		boxRect := Rect{w.Rect.X - 1, w.Rect.Y - 1, w.Rect.W + 2, w.Rect.H + 2}
		drawBox(boxRect, w.Fg, w.Bg)
	}

	starti := w.Scrollpos
	endi := w.Scrollpos + w.Rect.H - 1
	if endi > len(w.Items)-1 {
		endi = len(w.Items) - 1
	}

	y := w.Rect.Y
	for i := starti; i <= endi; i++ {
		item := w.Items[i]
		if w.Sel == i {
			// Highlight selected menu item
			if w.Settings&WidgetCenter != 0 {
				printpaddedcenter(item, w.Rect.X, y, w.Bg, w.Fg, w.Rect.W)
			} else {
				printpadded(item, 1, w.Rect.W-len(item)-1, w.Rect.X, y, w.Bg, w.Fg)
			}
		} else {
			if w.Settings&WidgetCenter != 0 {
				printpaddedcenter(item, w.Rect.X, y, w.Fg, w.Bg, w.Rect.W)
			} else {
				printpadded(item, 1, w.Rect.W-len(item)-1, w.Rect.X, y, w.Fg, w.Bg)
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
	starti := w.Scrollpos
	endi := w.Scrollpos + w.Rect.H - 1

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
