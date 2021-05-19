package main

import (
	tb "github.com/nsf/termbox-go"
)

type ListboxWidget struct {
	Rect      Rect
	Margin    Margin
	Color     Color
	Cb        WidgetEventCB
	Items     []*WidgetItem
	Settings  WidgetSetting
	Sel       int
	Scrollpos int
}

func NewListboxWidget(rect Rect, margin Margin, color Color, cb WidgetEventCB, items []*WidgetItem, settings WidgetSetting) *ListboxWidget {
	// If not specified, automatically set width and height based on listbox items.
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

	// Truncate listbox item text that go beyond width.
	for i, item := range items {
		if len(item.Display) > rect.W-margin.L-margin.R {
			items[i].Display = item.Display[:rect.W-margin.L-margin.R]
		}
	}

	InitColor(&color)

	w := ListboxWidget{
		Rect:      rect,
		Margin:    margin,
		Color:     color,
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
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&FmtBox != 0 {
		boxRect := AddRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
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
			// Highlight selected item.
			printspaces(w.Rect.W, w.Rect.X, y, w.Color.HighlightFg, w.Color.HighlightBg)
			printw(item.Display, contentRect.X, y, w.Color.HighlightFg, w.Color.HighlightBg, contentRect.W)
		} else {
			printw(item.Display, contentRect.X, y, w.Color.Fg, w.Color.Bg, contentRect.W)
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
				P1:   w.Items[w.Sel].Id,
				P2:   w.Sel,
				Pstr: w.Items[w.Sel].Sid,
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

func (w *ListboxWidget) SelItem() (int, *WidgetItem) {
	if len(w.Items) == 0 {
		return -1, nil
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
			P1:   w.Items[w.Sel].Id,
			P2:   w.Sel,
			Pstr: w.Items[w.Sel].Sid,
		}
		w.Cb(&we)
	}
}
