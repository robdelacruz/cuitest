package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	tb "github.com/nsf/termbox-go"
)

type WAccounts struct {
	db            *sql.DB
	Rect          TxRect
	Clr           TxColor
	Cb            TxEventCB
	tblAccounts   *TxTable
	tblSelAccount *TxTable
}

const (
	List int = iota
	ItemView
	ItemEdit
)

func NewWAccounts(db *sql.DB, rect TxRect, clr TxColor, cb TxEventCB) *WAccounts {
	initColor(&clr)

	w := WAccounts{
		db:   db,
		Rect: rect,
		Clr:  clr,
		Cb:   cb,
	}
	r := TxRect{0, 0, rect.W, rect.H}
	tblAccounts := createAccountsTable(db, r, clr, w.onAccountsEvent)

	w.tblAccounts = tblAccounts
	w.tblSelAccount = nil
	return &w
}

func createAccountsTable(db *sql.DB, r TxRect, clr TxColor, cb TxEventCB) *TxTable {
	props := &TxProps{r, TxMargin1, clr, cb, 0}
	cols := []*TxCellSetting{
		{"%s", 0, 40, clr, 0},
		{"%7.2f", 40, 12, clr, 0},
	}
	hh := []string{"Name", "Balance"}
	rows := queryAccountRows(db)
	return NewTxTable(props, clr, cols, hh, rows)
}

func queryAccountRows(db *sql.DB) []*TxTableRow {
	aa, err := findAccounts(db, " 1=1 ORDER BY accounttype, name")
	if err != nil {
		aa = []*Account{}
	}

	var rows []*TxTableRow
	for _, a := range aa {
		bal := balAccount(db, a.Accountid)
		cells := []TxCell{a.Name, bal}
		rows = append(rows, &TxTableRow{a.Accountid, a.Code, cells})
	}
	return rows
}

func (w *WAccounts) Draw() {
	clearRect(w.Rect, w.Clr.Bg)
	w.tblAccounts.Draw()
	tb.Flush()
}

func (w *WAccounts) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch == 0 {
		switch e.Key {
		case tb.KeyEnter: // view
			_log.Printf("view\n")
		}
		return w.tblAccounts.HandleEvent(e)
	}

	switch e.Ch {
	case 'a': // add
		_log.Printf("add\n")
	}
	return w.tblAccounts.HandleEvent(e)
}

func (w *WAccounts) onAccountsEvent(we *TxEvent) {
	switch we.Code {
	case TxEventEnter:
		_log.Printf("account enter, id: %d, alias: %s, display: %s\n", we.Item.Id, we.Item.Alias, we.Item.Display)
	case TxEventEsc:
	case TxEventSel:
	}
}
