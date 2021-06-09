package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type AccountType int

const (
	BankAccount AccountType = iota
	StockAccount
)

type Account struct {
	Accountid   int64       `json:"accountid"`
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	AccountType AccountType `json:"int"`
	Currencyid  int64       `json:"currencyid"`
}

func createAccount(db *sql.DB, a *Account) (int64, error) {
	s := "INSERT INTO account (code, name, accounttype, currency_id) VALUES (?, ?, ?, ?)"
	result, err := sqlexec(db, s, a.Code, a.Name, a.AccountType, a.Currencyid)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
func editAccount(db *sql.DB, a *Account) error {
	s := "UPDATE account SET code = ?, name = ?, accounttype = ?, currency_id = ? WHERE account_id = ?"
	_, err := sqlexec(db, s, a.Code, a.Name, a.AccountType, a.Currencyid, a.Accountid)
	if err != nil {
		return err
	}
	return nil
}
func delAccount(db *sql.DB, accountid int64) error {
	s := "DELETE FROM account WHERE account_id = ?"
	_, err := sqlexec(db, s, accountid)
	if err != nil {
		return err
	}
	return nil
}

func findAccount(db *sql.DB, accountid int64) (*Account, error) {
	s := "SELECT account_id, code, name, accounttype, currency_id FROM account WHERE account_id = ?"
	row := db.QueryRow(s, accountid)
	var a Account
	err := row.Scan(&a.Accountid, &a.Code, &a.Name, &a.AccountType, &a.Currencyid)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func findAccounts(db *sql.DB, swhere string) ([]*Account, error) {
	s := fmt.Sprintf("SELECT account_id, code, name, accounttype, currency_id FROM account WHERE %s", swhere)
	rows, err := db.Query(s)
	if err != nil {
		return nil, err
	}
	aa := []*Account{}
	for rows.Next() {
		var a Account
		rows.Scan(&a.Accountid, &a.Code, &a.Name, &a.AccountType, &a.Currencyid)
		aa = append(aa, &a)
	}
	return aa, nil
}

func balAccount(db *sql.DB, accountid int64) float64 {
	s := "SELECT IFNULL(SUM(amt), 0.0) FROM trans WHERE account_id = ?"
	row := db.QueryRow(s, accountid)
	var bal float64
	err := row.Scan(&bal)
	if err != nil {
		return 0.0
	}
	return bal
}
