package main

import (
	"fmt"
	tb "github.com/nsf/termbox-go"
)

type TxCellSetting struct {
	Sfmt string
	X, W int
	Clr  TxColor
	Fmt  TxFmt
}

type TxCell interface{}

type TxTableRow struct {
	Id    int64
	Alias string
	Cells []TxCell
}

type TxTable struct {
	Props      *TxProps
	HeadingClr TxColor
	Cols       []*TxCellSetting
	Headings   []string
	Rows       []*TxTableRow
	Sel        int
	Scrollpos  int
}

func NewTxTable(props *TxProps, headingclr TxColor, cols []*TxCellSetting, headings []string, rows []*TxTableRow) *TxTable {
	if props == nil {
		props = defaultProps()
	}

	// If not specified, automatically set width and height based on column settings.
	if props.Rect.H == 0 {
		props.Rect.H = len(rows) + props.Margin.T + props.Margin.B
		if len(headings) > 0 {
			props.Rect.H += 1
		}
	}
	if props.Rect.W == 0 {
		maxlen := 0
		for _, col := range cols {
			if col.X+col.W > maxlen {
				maxlen = col.X + col.W
			}
		}
		props.Rect.W = maxlen + props.Margin.L + props.Margin.R
	}

	initColor(&props.Clr)
	if headingclr.Fg == 0 {
		headingclr.Fg = props.Clr.Fg
	}
	if headingclr.Bg == 0 {
		headingclr.Bg = props.Clr.Bg
	}

	w := TxTable{
		Props:      props,
		HeadingClr: headingclr,
		Cols:       cols,
		Headings:   headings,
		Rows:       rows,
		Sel:        0,
		Scrollpos:  0,
	}

	return &w
}

func (w *TxTable) Draw() {
	p := w.Props
	clearRect(p.Rect, p.Clr.Bg)

	if p.Fmt&TxFmtBox != 0 {
		boxRect := addRectBox(p.Rect)
		drawBox(boxRect, p.Clr.Fg, p.Clr.Bg)
	}

	contentRect := addRectMargin(p.Rect, p.Margin)

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
			print(heading, contentRect.X+col.X, y, w.HeadingClr.Fg, w.HeadingClr.Bg)
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
			printspaces(p.Rect.W, p.Rect.X, y, p.Clr.HighlightFg, p.Clr.HighlightBg)
		}

		row := w.Rows[irow]
		for icol, cell := range row.Cells {
			if cell == nil {
				continue
			}
			if icol > len(w.Cols)-1 {
				continue
			}
			col := w.Cols[icol]

			// CellSetting color overrides Table color.
			fg := p.Clr.Fg
			if col.Clr.Fg != 0 {
				fg = col.Clr.Fg
			}
			bg := p.Clr.Bg
			if col.Clr.Bg != 0 {
				bg = col.Clr.Bg
			}

			// Use highlight color for cell when in selected row.
			if w.Sel == irow {
				fg = p.Clr.HighlightFg
				bg = p.Clr.HighlightBg
			}

			var scell string
			if col.Sfmt != "" {
				scell = fmt.Sprintf(col.Sfmt, cell)
			} else {
				scell = fmt.Sprintf("%v", cell)
			}
			if col.Fmt&TxFmtCenter != 0 {
				printcenter(scell, contentRect.X+col.X, y, fg, bg, col.W)
			} else {
				printw(scell, contentRect.X+col.X, y, fg, bg, col.W)
			}
		}
		y++
	}
}

func (w *TxTable) HandleEvent(e tb.Event) bool {
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
		w.adjustScroll()
		w.postSelItemEvent()
		return true
	case tb.KeyArrowDown:
		w.Sel++
		if w.Sel > len(w.Rows)-1 {
			w.Sel = 0
		}
		w.adjustScroll()
		w.postSelItemEvent()
		return true
	case tb.KeyEnter:
		if w.Props.EventCB == nil {
			return true
		}
		item := w.SelItem()
		if item == nil {
			return true
		}
		we := TxEvent{
			Code: TxEventEnter,
			Item: item,
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

func (w *TxTable) adjustScroll() {
	rect := addRectMargin(w.Props.Rect, w.Props.Margin)
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

func (w *TxTable) SelItem() *TxItem {
	if len(w.Rows) == 0 || w.Sel > len(w.Rows)-1 {
		return nil
	}
	row := w.Rows[w.Sel]
	item := &TxItem{row.Id, row.Alias, ""}
	return item
}

func (w *TxTable) postSelItemEvent() {
	if w.Props.EventCB == nil {
		return
	}
	item := w.SelItem()
	if item == nil {
		return
	}
	we := TxEvent{
		Code: TxEventSel,
		Item: item,
	}
	w.Props.EventCB(&we)
}

func (w *TxTable) SetRows(rows []*TxTableRow) {
	w.Rows = rows

	// Make sure selected row and scroll position is still in range.
	if w.Sel > len(w.Rows)-1 {
		w.Sel = len(w.Rows) - 1
	}
	if w.Scrollpos > len(w.Rows)-1 {
		w.Scrollpos = len(w.Rows) - 1
	}
}
