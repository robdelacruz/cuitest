package main

import (
	tb "github.com/nsf/termbox-go"
)

type TableRow []string
type CellSetting struct {
	X, W  int
	Color Color
}

type TableWidget struct {
	Rect         Rect
	Margin       Margin
	Color        Color
	ColorHeading Color
	Cb           WidgetEventCB
	Cols         []CellSetting
	Headings     []string
	Rows         []TableRow
	Settings     WidgetSetting
	Sel          int
	Scrollpos    int
}

func NewTableWidget(rect Rect, margin Margin, color Color, colorHeading Color, cb WidgetEventCB, cols []CellSetting, headings []string, rows []TableRow, settings WidgetSetting) *TableWidget {
	// If not specified, automatically set width and height based on column settings.
	if rect.H == 0 {
		rect.H = len(rows) + margin.T + margin.B
		if len(headings) > 0 {
			rect.H += 1
		}
	}
	if rect.W == 0 {
		maxlen := 0
		for _, col := range cols {
			if col.X+col.W > maxlen {
				maxlen = col.X + col.W
			}
		}
		rect.W = maxlen + margin.L + margin.R
	}

	InitColor(&color)
	if colorHeading.Fg == 0 {
		colorHeading.Fg = color.Fg
	}
	if colorHeading.Bg == 0 {
		colorHeading.Bg = color.Bg
	}

	w := TableWidget{
		Rect:         rect,
		Margin:       margin,
		Color:        color,
		ColorHeading: colorHeading,
		Cb:           cb,
		Cols:         cols,
		Headings:     headings,
		Rows:         rows,
		Settings:     settings,
		Sel:          0,
		Scrollpos:    0,
	}

	return &w
}

func (w *TableWidget) Draw() {
	clearRect(w.Rect, w.Color.Bg)

	if w.Settings&FmtBox != 0 {
		boxRect := AddRectBox(w.Rect)
		drawBox(boxRect, w.Color.Fg, w.Color.Bg)
	}

	contentRect := AddRectMargin(w.Rect, w.Margin)

	starti := w.Scrollpos
	endi := w.Scrollpos + contentRect.H - 1

	// Heading
	y := contentRect.Y
	if len(w.Headings) > 0 {
		for icol, heading := range w.Headings {
			if icol > len(w.Cols)-1 {
				continue
			}
			col := w.Cols[icol]
			print(heading, contentRect.X+col.X, y, w.ColorHeading.Fg, w.ColorHeading.Bg)
		}
		y++
		endi = endi - 1
	}

	if endi > len(w.Rows)-1 {
		endi = len(w.Rows) - 1
	}

	// Rows
	for irow := starti; irow <= endi; irow++ {
		if w.Sel == irow {
			// Highlight selected row.
			printspaces(w.Rect.W, w.Rect.X, y, w.Color.HighlightFg, w.Color.HighlightBg)
		}

		row := w.Rows[irow]
		for icol, cell := range row {
			if icol > len(w.Cols)-1 {
				continue
			}
			col := w.Cols[icol]
			fg := w.Color.Fg
			if col.Color.Fg != 0 {
				fg = col.Color.Fg
			}
			bg := w.Color.Bg
			if col.Color.Bg != 0 {
				bg = col.Color.Bg
			}
			if w.Sel == irow {
				// Use highlight color for cell when in selected row.
				printw(cell, contentRect.X+col.X, y, w.Color.HighlightFg, w.Color.HighlightBg, col.W)
			} else {
				printw(cell, contentRect.X+col.X, y, fg, bg, col.W)
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
	rect := AddRectMargin(w.Rect, w.Margin)
	starti := w.Scrollpos
	endi := w.Scrollpos + rect.H - 1

	if len(w.Headings) > 0 {
		endi = endi - 1
	}

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
