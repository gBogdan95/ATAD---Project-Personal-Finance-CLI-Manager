package db

import (
	"database/sql"
	"fmt"

	"github.com/gBogdan95/ATAD---Project-Personal-Finance-CLI-Manager/internal/models"
)

func ListTransactions(conn *sql.DB, limit int) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := conn.Query(
		`SELECT id, date, type, amount_cents, category, COALESCE(note, '')
		 FROM transactions
		 ORDER BY date DESC, id DESC
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var out []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Type, &t.AmountCents, &t.Category, &t.Note); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}

	return out, nil
}
