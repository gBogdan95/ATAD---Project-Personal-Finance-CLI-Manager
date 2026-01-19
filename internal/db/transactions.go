package db

import (
	"database/sql"
	"fmt"

	"github.com/gBogdan95/ATAD---Project-Personal-Finance-CLI-Manager/internal/models"
)

func InsertTransaction(conn *sql.DB, t models.Transaction) (int64, error) {
	res, err := conn.Exec(
		`INSERT INTO transactions(date, type, amount_cents, category, note)
		 VALUES(?, ?, ?, ?, ?)`,
		t.Date, t.Type, t.AmountCents, t.Category, t.Note,
	)
	if err != nil {
		return 0, fmt.Errorf("insert transaction: %w", err)
	}
	return res.LastInsertId()
}
