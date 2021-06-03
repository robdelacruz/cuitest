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

	w := WAccounts{
		db:   db,
		rect: rect,
		tbl:  nil,
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
