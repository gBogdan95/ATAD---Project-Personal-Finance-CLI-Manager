package db

import (
	"database/sql"
	"fmt"
)

type MonthSummary struct {
	IncomeCents  int64
	ExpenseCents int64
}

func GetMonthSummary(conn *sql.DB, month string) (MonthSummary, error) {
	prefix := month + "-"

	var income int64
	if err := conn.QueryRow(
		`SELECT COALESCE(SUM(amount_cents), 0)
		 FROM transactions
		 WHERE type = 'income' AND date LIKE ?`,
		prefix+"%",
	).Scan(&income); err != nil {
		return MonthSummary{}, fmt.Errorf("sum income: %w", err)
	}

	var expense int64
	if err := conn.QueryRow(
		`SELECT COALESCE(SUM(amount_cents), 0)
		 FROM transactions
		 WHERE type = 'expense' AND date LIKE ?`,
		prefix+"%",
	).Scan(&expense); err != nil {
		return MonthSummary{}, fmt.Errorf("sum expense: %w", err)
	}

	return MonthSummary{
		IncomeCents:  income,
		ExpenseCents: expense,
	}, nil
}
