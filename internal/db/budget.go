package db

import (
	"database/sql"
	"fmt"

	"github.com/gBogdan95/ATAD---Project-Personal-Finance-CLI-Manager/internal/models"
)

func UpsertBudget(conn *sql.DB, month, category string, limitCents int64) error {
	_, err := conn.Exec(
		`INSERT INTO budgets(month, category, limit_cents)
		 VALUES(?, ?, ?)
		 ON CONFLICT(month, category) DO UPDATE SET limit_cents = excluded.limit_cents`,
		month, category, limitCents,
	)
	if err != nil {
		return fmt.Errorf("upsert budget: %w", err)
	}
	return nil
}

func ListBudgets(conn *sql.DB, month string) ([]models.Budget, error) {
	rows, err := conn.Query(
		`SELECT id, month, category, limit_cents
		 FROM budgets
		 WHERE month = ?
		 ORDER BY category ASC`,
		month,
	)
	if err != nil {
		return nil, fmt.Errorf("query budgets: %w", err)
	}
	defer rows.Close()

	var out []models.Budget
	for rows.Next() {
		var b models.Budget
		if err := rows.Scan(&b.ID, &b.Month, &b.Category, &b.LimitCents); err != nil {
			return nil, fmt.Errorf("scan budget: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate budgets: %w", err)
	}
	return out, nil
}

func SumExpensesForCategoryMonth(conn *sql.DB, month, category string) (int64, error) {
	prefix := month + "-"

	var spent int64
	if err := conn.QueryRow(
		`SELECT COALESCE(SUM(amount_cents), 0)
		 FROM transactions
		 WHERE type = 'expense'
		   AND category = ?
		   AND date LIKE ?`,
		category, prefix+"%",
	).Scan(&spent); err != nil {
		return 0, fmt.Errorf("sum expenses: %w", err)
	}

	return spent, nil
}
