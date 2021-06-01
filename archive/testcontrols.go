package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	tb "github.com/nsf/termbox-go"
)

var _log *log.Logger
var _termW, _termH int

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func run(args []string) error {
	flog, err := os.Create("./log.txt")
	if err != nil {
		return err
	}
	defer flog.Close()
	_log = log.New(flog, "", 0)

	sw, parms := parseArgs(args)

	// [-i new_file]  Create and initialize db file
	if sw["i"] != "" {
		dbfile := sw["i"]
		if fileExists(dbfile) {
			return fmt.Errorf("File '%s' already exists. Can't initialize it.\n", dbfile)
		}
		createTables(dbfile)
		return nil
	}

	// Need to specify a db file as first parameter.
	if len(parms) == 0 {
		s := `Usage:

   Specify database file:
	t <db file>

   To initialize new database file:
	t -i <new db file>

`
		fmt.Printf(s)
		return nil
	}

	// Exit if db file doesn't exist.
	dbfile := parms[0]
	if !fileExists(dbfile) {
		return fmt.Errorf(`Database file '%s' doesn't exist. Create one using:
	t -i <filename>
   `, dbfile)
	}

	//db, err := sql.Open("sqlite3", dbfile)
	_, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		return fmt.Errorf("Error opening '%s' (%s)\n", dbfile, err)
	}

	// Start termbox mode
	err = tb.Init()
	if err != nil {
		return err
	}
	defer tb.Close()
	tb.SetOutputMode(tb.Output256)

	_termW, _termH = tb.Size()
	mainwin := NewMainWindow()
	mainwin.Draw()

	chev := make(chan tb.Event)
	go pollEvents(chev)

	for {
		e := <-chev

		if e.Ch == 'q' {
			break
		}
		if mainwin.HandleEvent(e) {
			mainwin.Draw()
		}
	}
	return nil
}

func listContains(ss []string, v string) bool {
	for _, s := range ss {
		if v == s {
			return true
		}
	}
	return false
}
func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func pollEvents(chev chan tb.Event) {
	for {
		e := tb.PollEvent()
		if e.Type != tb.EventKey {
			continue
		}
		chev <- e
	}
}

func createTables(newfile string) {
	if fileExists(newfile) {
		s := fmt.Sprintf("File '%s' already exists. Can't initialize it.\n", newfile)
		fmt.Printf(s)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", newfile)
	if err != nil {
		fmt.Printf("Error opening '%s' (%s)\n", newfile, err)
		os.Exit(1)
	}

	ss := []string{
		"CREATE TABLE currency (currency_id INTEGER PRIMARY KEY NOT NULL, name TEXT, usdrate REAL);",
		"CREATE TABLE account (account_id INTEGER PRIMARY KEY NOT NULL, code TEXT, name TEXT, accounttype INTEGER, currency_id INTEGER);",
		"CREATE TABLE trans (trans_id INTEGER PRIMARY KEY NOT NULL, account_id INTEGER, date TEXT, ref TEXT, desc TEXT, amt REAL);",
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("DB error (%s)\n", err)
		os.Exit(1)
	}
	for _, s := range ss {
		_, err := txexec(tx, s)
		if err != nil {
			tx.Rollback()
			log.Printf("DB error (%s)\n", err)
			os.Exit(1)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("DB error (%s)\n", err)
		os.Exit(1)
	}

	initTestData(db)
}

func initTestData(db *sql.DB) {
	c1 := Currency{
		Name:    "USD",
		Usdrate: 1.0,
	}
	c2 := Currency{
		Name:    "PHP",
		Usdrate: 48.0,
	}
	usdid, err := createCurrency(db, &c1)
	if err != nil {
		panic(err)
	}
	phpid, err := createCurrency(db, &c2)
	if err != nil {
		panic(err)
	}

	a1 := Account{
		Code:        "bpichecking",
		Name:        "BPI Checking Account",
		AccountType: BankAccount,
		Currencyid:  phpid,
	}
	a2 := Account{
		Code:        "bpisavings",
		Name:        "BPI Savings Account",
		AccountType: BankAccount,
		Currencyid:  phpid,
	}
	a3 := Account{
		Code:        "bpiusd",
		Name:        "BPI USD",
		AccountType: BankAccount,
		Currencyid:  usdid,
	}
	_, err = createAccount(db, &a1)
	if err != nil {
		panic(err)
	}
	_, err = createAccount(db, &a2)
	if err != nil {
		panic(err)
	}
	_, err = createAccount(db, &a3)
	if err != nil {
		panic(err)
	}
}

func parseArgs(args []string) (map[string]string, []string) {
	switches := map[string]string{}
	parms := []string{}

	standaloneSwitches := []string{}
	definitionSwitches := []string{"i"}
	fNoMoreSwitches := false
	curKey := ""

	for _, arg := range args {
		if fNoMoreSwitches {
			// any arg after "--" is a standalone parameter
			parms = append(parms, arg)
		} else if arg == "--" {
			// "--" means no more switches to come
			fNoMoreSwitches = true
		} else if strings.HasPrefix(arg, "--") {
			switches[arg[2:]] = "y"
			curKey = ""
		} else if strings.HasPrefix(arg, "-") {
			if listContains(definitionSwitches, arg[1:]) {
				// -a "val"
				curKey = arg[1:]
				continue
			}
			for _, ch := range arg[1:] {
				// -a, -b, -ab
				sch := string(ch)
				if listContains(standaloneSwitches, sch) {
					switches[sch] = "y"
				}
			}
		} else if curKey != "" {
			switches[curKey] = arg
			curKey = ""
		} else {
			// standalone parameter
			parms = append(parms, arg)
		}
	}

	return switches, parms
}

type MainWindow struct {
	width, height int

	smileyw *Smiley
	labelw  *TxLabel
	popupw  TxWidget
}

func NewMainWindow() *MainWindow {
	w := MainWindow{}
	w.width = _termW
	w.height = _termH

	w.smileyw = &Smiley{}
	w.smileyw.X = 5
	w.smileyw.Y = 5

	grey39 := tb.Attribute(242)
	gold1 := tb.Attribute(221)

	color := TxColor{
		Fg: gold1,
		Bg: grey39,
	}

	lbltext := "Now is the time for all good men to come to the aid of the party.\n\nThe quick brown fox jumps over the lazy dog. Now is the time for all good men to come to the aid of the party. The quick brown fox jumps over the lazy dog. Now is the time for all good men to come to the aid of the party. The quick brown fox jumps over the lazy dog."

	w.labelw = NewTxLabel(TxRect{1, 20, 30, 10}, TxMarginX, color, lbltext, TxFmtBox)

	return &w
}

func (w *MainWindow) popupCB(we *TxEvent) {
	if we.Code == TxEventEnter {
		_, ok := w.popupw.(*TxMenu)
		if ok {
			sel := we.P1
			item := we.Item
			w.labelw.SetText(fmt.Sprintf("Selected menu option %d: %s\n", sel, item.Display))
		}
		_, ok = w.popupw.(*TxListbox)
		if ok {
			w.labelw.SetText(fmt.Sprintf("Selected listbox item:\n%s", we.Item.Display))
		}
		w.popupw = nil
	} else if we.Code == TxEventEsc {
		w.labelw.SetText("*** Canceled operation ***")
		w.popupw = nil
	} else if we.Code == TxEventSel {
		if we.Item != nil {
			w.labelw.SetText(fmt.Sprintf("[%s]: %s", we.Item.Sid, we.Item.Display))
		}
	}
}

func (w *MainWindow) Draw() {
	tb.Clear(0, 0)

	w.smileyw.Draw()
	w.labelw.Draw()

	if w.popupw != nil {
		w.popupw.Draw()
	}

	tb.Flush()
}

func (w *MainWindow) HandleEvent(e tb.Event) bool {
	if w.popupw != nil {
		return w.popupw.HandleEvent(e)
	}

	white := tb.Attribute(16)
	black := tb.Attribute(17)
	darkolivegreen := tb.Attribute(156)
	darkorange := tb.Attribute(167)
	grey39 := tb.Attribute(242)
	plum1 := tb.Attribute(220)
	gold1 := tb.Attribute(221)

	color1 := TxColor{
		Fg:          darkolivegreen,
		Bg:          black,
		HighlightFg: grey39,
		HighlightBg: white,
	}
	color2 := TxColor{
		Fg:          darkorange,
		Bg:          black,
		HighlightFg: gold1,
		HighlightBg: grey39,
	}
	color3 := TxColor{
		Fg: plum1,
		Bg: grey39,
	}
	color4 := TxColor{
		Fg: darkolivegreen,
		Bg: grey39,
	}

	if e.Ch == 'm' {
		items := []*TxItem{
			&TxItem{1, "option1", "Menu Option 1 abc"},
			&TxItem{2, "option2", "Option 2 def"},
			&TxItem{3, "option3", "Option 3 ghijkl"},
			&TxItem{4, "option4", "Option 4 some more text"},
			&TxItem{5, "option5", "Option 5 xyz"},
			&TxItem{6, "option6", "Option 6 lmnop"},
			&TxItem{7, "option7", "Option 7 qrstuvw"},
			&TxItem{8, "option8", "Option 8 12345"},
			&TxItem{9, "option9", "Option 9 123"},
			&TxItem{10, "option10", "Option 10"},
		}
		//w.popupw = NewTxMenu(TxRect{5, 1, 0, 0}, TxMargin0, color1, w.popupCB, items, TxFmtBox|TxFmtCenter)
		w.popupw = NewTxMenu(TxRect{5, 1, 31, 0}, TxMarginX, color1, w.popupCB, items, TxFmtBox|TxFmtCenter)
		return true
	} else if e.Ch == 'l' {
		items := []*TxItem{
			&TxItem{1, "line1", "Now is the time"},
			&TxItem{2, "line2", "for all good men"},
			&TxItem{3, "line3", "to come to the aid"},
			&TxItem{4, "line4", "of the party."},
			&TxItem{5, "line5", "-- typing drill"},
		}
		w.popupw = NewTxListbox(TxRect{10, 1, 30, 0}, TxMarginX, color2, w.popupCB, items, TxFmtBox)
		return true
	} else if e.Ch == 't' {
		var cellColor TxColor
		cols := []CellSetting{
			CellSetting{0, 15, cellColor},
			CellSetting{15, 10, cellColor},
			CellSetting{35, 10, cellColor},
		}
		headings := []string{"col1", "col2", "col3"}
		rows := []TableRow{
			TableRow{"abc", "defghi", "jklmn"},
			TableRow{"ABC", "DEFGHI", "JKLMN"},
			TableRow{"12345", "678", "9012"},
			TableRow{"Now is", "the time", "for all"},
			TableRow{"good men", "to come to", "the aid"},
			TableRow{"of the", "party.", ""},
			TableRow{"12345", "678", "9012"},
		}
		//w.popupw = NewTxTable(TxRect{5, 5, 0, 0}, color3, color1, w.popupCB, cols, headings, rows, TxFmtBox)
		w.popupw = NewTxTable(TxRect{5, 5, 0, 0}, TxMargin1, color3, color1, w.popupCB, cols, headings, rows, TxFmtCenter)
		return true
	} else if e.Ch == 'e' {
		w.popupw = NewTxEntry(TxRect{5, 5, 30, 1}, TxMargin0, color4, w.popupCB, "Entry Text", TxFmtBox)
		return true
	}

	switch e.Key {
	}

	return w.smileyw.HandleEvent(e)
}

type Smiley struct {
	X, Y int
}

func (sm *Smiley) Draw() {
	fg := tb.Attribute(156)
	bg := tb.Attribute(0)
	tb.SetCell(sm.X, sm.Y, '😀', fg, bg)
}

func (sm *Smiley) HandleEvent(e tb.Event) bool {
	if e.Type != tb.EventKey {
		return false
	}
	if e.Ch != 0 {
		return false
	}

	switch e.Key {
	case tb.KeyArrowUp:
		if sm.Y > 0 {
			sm.Y--
			return true
		}
	case tb.KeyArrowDown:
		if sm.Y < _termH-1 {
			sm.Y++
			return true
		}
	case tb.KeyArrowLeft:
		if sm.X > 0 {
			sm.X--
			return true
		}
	case tb.KeyArrowRight:
		if sm.X < _termW-1 {
			sm.X++
			return true
		}
	}
	return false
}
