package main

import (
	tb "github.com/nsf/termbox-go"
)

type TableRow []string
type CellSetting struct {
	X, W   int
	Fg, Bg tb.Attribute
}

type TableWidget struct {
	Rect                 Rect
	Fg, Bg               tb.Attribute
	HeadingFg, HeadingBg tb.Attribute
	Cb                   WidgetEventCB
	Cols                 []CellSetting
	Headings             []string
	Rows                 []TableRow
	Settings             WidgetSetting
	Sel                  int
	Scrollpos            int
}

func NewTableWidget(rect Rect, fg, bg, headingfg, headingbg tb.Attribute, cb WidgetEventCB, cols []CellSetting, headings []string, rows []TableRow, settings WidgetSetting) *TableWidget {
	// If not specified, automatically set width and height based on column settings.
	if rect.H == 0 {
		rect.H = len(rows)
	}
	if rect.W == 0 {
		maxlen := 0
		for _, col := range cols {
			if col.X+col.W > maxlen {
				maxlen = col.X + col.W
			}
		}

		// Add 1 char margin to the right.
		rect.W = maxlen + 1
	}

	w := TableWidget{
		Rect:      rect,
		Fg:        fg,
		Bg:        bg,
		HeadingFg: headingfg,
		HeadingBg: headingbg,
		Cb:        cb,
		Cols:      cols,
		Headings:  headings,
		Rows:      rows,
		Settings:  settings,
		Sel:       0,
		Scrollpos: 0,
	}

	return &w
}

func (w *TableWidget) Draw() {
	clearRect(w.Rect, w.Bg)

	if w.Settings&WidgetBox != 0 {
		boxRect := Rect{w.Rect.X - 1, w.Rect.Y - 1, w.Rect.W + 2, w.Rect.H + 2}
		drawBox(boxRect, w.Fg, w.Bg)
	}

	starti := w.Scrollpos
	endi := w.Scrollpos + w.Rect.H - 1
	if endi > len(w.Rows)-1 {
		endi = len(w.Rows) - 1
	}

	y := w.Rect.Y
	for irow := starti; irow <= endi; irow++ {
		if w.Sel == irow {
			printspaces(w.Rect.W, w.Rect.X, y, w.Bg, w.Fg)
		}

		row := w.Rows[irow]
		for icol, cell := range row {
			if icol > len(w.Cols)-1 {
				continue
			}
			col := w.Cols[icol]
			fg := w.Fg
			if col.Fg != 0 {
				fg = col.Fg
			}
			bg := w.Bg
			if col.Bg != 0 {
				bg = col.Bg
			}
			if w.Sel == irow {
				// Highlight selected row
				print(cell, w.Rect.X+col.X, y, bg, fg)
			} else {
				print(cell, w.Rect.X+col.X, y, fg, bg)
			}
		}
		y++
	}
}

func (w *TableWidget) HandleEvent(e tb.Event) bool {
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
			w.Sel = len(w.Rows) - 1
		}
		w.AdjustScroll()
		w.PostSelRowEvent()
		return true
	case tb.KeyArrowDown:
		w.Sel++
		if w.Sel > len(w.Rows)-1 {
			w.Sel = 0
		}
		w.AdjustScroll()
		w.PostSelRowEvent()
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

func (w *TableWidget) AdjustScroll() {
	starti := w.Scrollpos
	endi := w.Scrollpos + w.Rect.H - 1

	if w.Sel < starti {
		w.Scrollpos -= starti - w.Sel
	} else if w.Sel > endi {
		w.Scrollpos += w.Sel - endi
	}

	if w.Scrollpos < 0 {
		w.Scrollpos = 0
	} else if w.Scrollpos > len(w.Rows)-1 {
		w.Scrollpos = len(w.Rows) - 1
	}
}

func (w *TableWidget) SelItem() (int, TableRow) {
	if len(w.Rows) == 0 {
		return -1, nil
	}
	return w.Sel, w.Rows[w.Sel]
}

func (w *TableWidget) PostSelRowEvent() {
	if len(w.Rows) == 0 {
		return
	}

	if w.Cb != nil {
		we := WidgetEvent{
			Code: WidgetEventSel,
			P1:   w.Sel,
		}
		w.Cb(&we)
	}
}
