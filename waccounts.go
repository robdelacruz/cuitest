package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	tb "github.com/nsf/termbox-go"
)

type WAccounts struct {
	db   *sql.DB
	rect TxRect
	clr  TxColor
	tbl  *TxTable
}

func NewWAccounts(db *sql.DB, rect TxRect, clr TxColor) *WAccounts {
	initColor(&clr)

	cols := []CellSetting{
		CellSetting{"%s", 0, 40, clr},
		CellSetting{"%s", 40, 12, clr},
		CellSetting{"%s", 52, 12, clr},
	}
	headings := []string{"Desc", "Deposit", "Withdraw"}
	rows := [][]interface{}{
		[]interface{}{"Initial deposit", "12345.67", ""},
		[]interface{}{"withdraw", "", "100.00"},
		[]interface{}{"transfer to savings", "", "550.25"},
		[]interface{}{"deposit paycheck", "345.67", ""},
		[]interface{}{"interest", "2.80", ""},
	}
	tbl := NewTxTable(rect, TxMargin0, clr, TxColor{}, nil, cols, headings, rows, 0)
	w := WAccounts{
		db:   db,
		rect: rect,
		tbl:  tbl,
	}
	return &w
}

func (w *WAccounts) Draw() {
	clearRect(w.rect, w.clr.Bg)
	w.tbl.Draw()
	tb.Flush()
}

func (w *WAccounts) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch != 0 {
		return false
	}
	return w.tbl.HandleEvent(e)
}
