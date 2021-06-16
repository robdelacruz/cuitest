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

	//white := tb.Attribute(16)
	//black := tb.Attribute(17)
	//darkolivegreen := tb.Attribute(156)
	//darkorange := tb.Attribute(167)
	//grey39 := tb.Attribute(242)
	//plum1 := tb.Attribute(220)
	//gold1 := tb.Attribute(221)

	//r := TxRect{0, 0, 80, 25}
	//waccounts := NewWAccounts(db, r, TxColorBW)
	//waccounts.Draw()

	r := TxRect{5, 5, 40, 1}
	entry := NewTxEntry(r, TxMargin0, TxColorBW, nil, "", 0)
	entry.Draw()

	tb.Flush()

	chev := make(chan tb.Event)
	go pollEvents(chev)

	for {
		e := <-chev

		if e.Ch == 'q' {
			break
		}
		if entry.HandleEvent(e) {
			entry.Draw()
			tb.Flush()
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
