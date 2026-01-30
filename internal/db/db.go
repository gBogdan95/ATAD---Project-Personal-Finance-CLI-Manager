package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const appDirName = "pfm"
const dbFileName = "pfm.db"

func DBPath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(cfgDir, appDirName, dbFileName), nil
}

func EnsureDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	dir := filepath.Join(cfgDir, appDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}
	return dir, nil
}

func Open(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)"
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return conn, nil
}

func InitSchema(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS transactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date TEXT NOT NULL,
  type TEXT NOT NULL CHECK(type IN ('income','expense')),
  amount_cents INTEGER NOT NULL,
  category TEXT NOT NULL,
  note TEXT
);

CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);

CREATE TABLE IF NOT EXISTS budgets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  month TEXT NOT NULL,
  category TEXT NOT NULL,
  limit_cents INTEGER NOT NULL,
  UNIQUE(month, category)
);

CREATE INDEX IF NOT EXISTS idx_budgets_month ON budgets(month);
`
	_, err := conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}
