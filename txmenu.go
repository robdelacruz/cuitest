package main

import (
	tb "github.com/nsf/termbox-go"
)

type TxMenu struct {
	Props     *TxProps
	Items     []*TxItem
	Sel       int
	Scrollpos int
}

func NewTxMenu(props *TxProps, items []*TxItem) *TxMenu {
	if props == nil {
		props = defaultProps()
	}

	// If not specified, automatically set width and height based on menu items.
	if props.Rect.H == 0 {
		props.Rect.H = len(items) + props.Margin.T + props.Margin.B
	}
	if props.Rect.W == 0 {
		maxlen := 0
		for _, item := range items {
			if len(item.Display) > maxlen {
				maxlen = len(item.Display)
			}
		}
		props.Rect.W = maxlen + props.Margin.L + props.Margin.R
	}

	// Truncate menu item text that go beyond width.
	for i, item := range items {
		if len(item.Display) > props.Rect.W-props.Margin.L-props.Margin.R {
			items[i].Display = item.Display[:props.Rect.W-props.Margin.L-props.Margin.R]
		}
	}

	initColor(&props.Clr)

	w := TxMenu{
		Props:     props,
		Items:     items,
		Sel:       0,
		Scrollpos: 0,
	}
	return &w
}

func (w *TxMenu) Draw() {
	p := w.Props
	clearRect(p.Rect, p.Clr.Bg)

	if p.Fmt&TxFmtBox != 0 {
		boxRect := addRectBox(p.Rect)
		drawBox(boxRect, p.Clr.Fg, p.Clr.Bg)
	}

	contentRect := addRectMargin(p.Rect, p.Margin)

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
			printspaces(p.Rect.W, p.Rect.X, y, p.Clr.HighlightFg, p.Clr.HighlightBg)
			if p.Fmt&TxFmtCenter != 0 {
				printcenter(item.Display, contentRect.X, y, p.Clr.HighlightFg, p.Clr.HighlightBg, contentRect.W)
				printcenter(item.Display, p.Rect.X, y, p.Clr.HighlightFg, p.Clr.HighlightBg, p.Rect.W)
			} else {
				printw(item.Display, contentRect.X, y, p.Clr.HighlightFg, p.Clr.HighlightBg, contentRect.W)
			}
		} else {
			if p.Fmt&TxFmtCenter != 0 {
				printcenter(item.Display, contentRect.X, y, p.Clr.Fg, p.Clr.Bg, contentRect.W)
			} else {
				printw(item.Display, contentRect.X, y, p.Clr.Fg, p.Clr.Bg, contentRect.W)
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
		if w.Props.EventCB == nil {
			return true
		}
		if len(w.Items) == 0 || w.Sel > len(w.Items)-1 {
			return true
		}
		we := TxEvent{
			Code: TxEventEnter,
			Item: w.Items[w.Sel],
		}
		w.Props.EventCB(&we)
		return true
	case tb.KeyEsc:
		if w.Props.EventCB == nil {
			return true
		}
		we := TxEvent{
			Code: TxEventEsc,
		}
		w.Props.EventCB(&we)
		return true
	}
	return false
}

func (w *TxMenu) adjustScroll() {
	rect := addRectMargin(w.Props.Rect, w.Props.Margin)
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
