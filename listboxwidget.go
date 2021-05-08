package main

import (
	tb "github.com/nsf/termbox-go"
)

type ListboxWidget struct {
	Rect      Rect
	Attrs     WidgetAttributes
	Cb        WidgetEventCB
	Items     []string
	Settings  WidgetSetting
	Sel       int
	Scrollpos int
}

func NewListboxWidget(rect Rect, attrs WidgetAttributes, cb WidgetEventCB, items []string, settings WidgetSetting) *ListboxWidget {
	// If not specified, automatically set width and height based on listbox items.
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

	// Truncate listbox item text that go beyond width.
	for i, item := range items {
		if len(item)+2 > rect.W {
			items[i] = item[:rect.W-2]
		}
	}

	InitWidgetAttributes(&attrs)

	w := ListboxWidget{
		Rect:      rect,
		Attrs:     attrs,
		Cb:        cb,
		Items:     items,
		Settings:  settings,
		Sel:       0,
		Scrollpos: 0,
	}
	w.PostSelItemEvent()
	return &w
}

func (w *ListboxWidget) Draw() {
	clearRect(w.Rect, w.Attrs.Bg)

	if w.Settings&WidgetBox != 0 {
		boxRect := Rect{w.Rect.X - 1, w.Rect.Y - 1, w.Rect.W + 2, w.Rect.H + 2}
		drawBox(boxRect, w.Attrs.Fg, w.Attrs.Bg)
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
			// Highlight selected item
			printpadded(item, 1, w.Rect.W-len(item)-1, w.Rect.X, y, w.Attrs.HighlightFg, w.Attrs.HighlightBg)
		} else {
			printpadded(item, 1, w.Rect.W-len(item)-1, w.Rect.X, y, w.Attrs.Fg, w.Attrs.Bg)
		}
		y++
	}
}

func (w *ListboxWidget) HandleEvent(e tb.Event) bool {
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
		w.PostSelItemEvent()
		return true
	case tb.KeyArrowDown:
		w.Sel++
		if w.Sel > len(w.Items)-1 {
			w.Sel = 0
		}
		w.AdjustScroll()
		w.PostSelItemEvent()
		return true
	case tb.KeyEnter:
		if w.Cb != nil {
			we := WidgetEvent{
				Code: WidgetEventEnter,
				P1:   w.Sel,
				Pstr: w.Items[w.Sel],
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

func (w *ListboxWidget) AdjustScroll() {
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

func (w *ListboxWidget) SelItem() (int, string) {
	if len(w.Items) == 0 {
		return -1, ""
	}
	return w.Sel, w.Items[w.Sel]
}

func (w *ListboxWidget) PostSelItemEvent() {
	if len(w.Items) == 0 {
		return
	}

	if w.Cb != nil {
		we := WidgetEvent{
			Code: WidgetEventSel,
			P1:   w.Sel,
			Pstr: w.Items[w.Sel],
		}
		w.Cb(&we)
	}
}
