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

	aa, err := findAccounts(db, " 1=1 ORDER BY accounttype, name")
	if err != nil {
		aa = []*Account{}
	}

	cols := []TxCellSetting{
		{"%s", 0, 40, TxColorBW, 0},
		{"%12.2f", 40, 12, TxColorBW, 0},
	}
	hh := []string{"Name", "Balance"}
	var rows []TxTableRow
	for _, a := range aa {
		bal := balAccount(db, a.Accountid)
		rows = append(rows, TxTableRow{a.Name, bal})
	}
	r := TxRect{0, 0, rect.W, rect.H}
	tbl := NewTxTable(r, TxMargin1, clr, clr, nil, cols, hh, rows, 0)

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
