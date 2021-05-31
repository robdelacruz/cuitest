package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

//CREATE TABLE currency (currency_id INTEGER PRIMARY KEY NOT NULL, name TEXT, usdrate REAL)

type Currency struct {
	Currencyid int64   `json:"currencyid"`
	Name       string  `json:"name"`
	Usdrate    float64 `json:"usdrate"`
}

func createCurrency(db *sql.DB, c *Currency) (int64, error) {
	s := "INSERT INTO currency (name, usdrate) VALUES (?, ?)"
	result, err := sqlexec(db, s, c.Name, c.Usdrate)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func editCurrency(db *sql.DB, c *Currency) error {
	s := "UPDATE currency SET name = ?, usdrate = ? WHERE currency_id = ?"
	_, err := sqlexec(db, s, c.Name, c.Usdrate, c.Currencyid)
	if err != nil {
		return err
	}
	return nil
}
func delCurrency(db *sql.DB, currencyid int64) error {
	s := "DELETE FROM currency WHERE currency_id = ?"
	_, err := sqlexec(db, s, currencyid)
	if err != nil {
		return err
	}
	return nil
}

func findCurrency(db *sql.DB, currencyid int64) (*Currency, error) {
	s := "SELECT currency_id, name, usdrate, FROM currency WHERE currency_id = ?"
	row := db.QueryRow(s, currencyid)
	var c Currency
	err := row.Scan(&c.Currencyid, &c.Name, &c.Usdrate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func findCurrencies(db *sql.DB, swhere string) ([]*Currency, error) {
	s := fmt.Sprintf("SELECT currency_id, name, usdrate FROM currency WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	cc := []*Currency{}
	for rows.Next() {
		var c Currency
		rows.Scan(&c.Currencyid, &c.Name, &c.Usdrate)
		cc = append(cc, &c)
	}
	return cc, nil
}
