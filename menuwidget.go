package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxMenu struct {
	Rect      TxRect
	Margin    TxMargin
	Color     TxColor
	Cb        TxEventCB
	Items     []*TxItem
	Settings  TxSetting
	Sel       int
	Scrollpos int
}

func NewTxMenu(rect TxRect, margin TxMargin, color TxColor, cb TxEventCB, items []*TxItem, settings TxSetting) *TxMenu {
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

	initColor(&color)

	w := TxMenu{
		Rect:      rect,
		Margin:    margin,
		Color:     color,
		Cb:        cb,
		Items:     items,
		Settings:  settings,
		Sel:       0,
		Scrollpos: 0,
	}
	return &w
}

func (w *TxMenu) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&TxFmtBox != 0 {
		boxRect := addRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	contentRect := addRectMargin(w.Rect, w.Margin)

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
			printspaces(w.Rect.W, w.Rect.X, y, w.Color.HighlightFg, w.Color.HighlightBg)
			if w.Settings&TxFmtCenter != 0 {
				printcenter(item.Display, contentRect.X, y, w.Color.HighlightFg, w.Color.HighlightBg, contentRect.W)
				printcenter(item.Display, w.Rect.X, y, w.Color.HighlightFg, w.Color.HighlightBg, w.Rect.W)
			} else {
				printw(item.Display, contentRect.X, y, w.Color.HighlightFg, w.Color.HighlightBg, contentRect.W)
			}
		} else {
			if w.Settings&TxFmtCenter != 0 {
				printcenter(item.Display, contentRect.X, y, w.Color.Fg, w.Color.Bg, contentRect.W)
			} else {
				printw(item.Display, contentRect.X, y, w.Color.Fg, w.Color.Bg, contentRect.W)
			}
		}
		y++
	}
}

func (w *TxMenu) HandleEvent(e tb.Event) bool {
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
		w.adjustScroll()
		return true
	case tb.KeyArrowDown:
		w.Sel++
		if w.Sel > len(w.Items)-1 {
			w.Sel = 0
		}
		w.adjustScroll()
		return true
	case tb.KeyEnter:
		if w.Cb != nil {
			we := TxEvent{
				Code: TxEventEnter,
				P1:   w.Sel,
				Item: w.Items[w.Sel],
			}
			w.Cb(&we)
		}
		return true
	case tb.KeyEsc:
		if w.Cb != nil {
			we := TxEvent{
				Code: TxEventEsc,
			}
			w.Cb(&we)
		}
		return true
	}
	return false
}

func (w *TxMenu) adjustScroll() {
	rect := addRectMargin(w.Rect, w.Margin)
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
