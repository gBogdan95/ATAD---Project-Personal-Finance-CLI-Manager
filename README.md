# Personal Finance CLI Manager (pfm)

A command-line app for tracking income and expenses locally using SQLite.

## Milestone 1 (Architecture + Basic Functionality)
- Go project structure created
- CLI commands implemented: `hello`, `version`, `init`, `add`, `list`
- `init` creates a local SQLite database and the initial `transactions` table
- `add` inserts a transaction into the database
- `list` displays the most recent transactions (default 20, configurable via `--limit`)

## Run
```bash
go run ./cmd/pfm hello
go run ./cmd/pfm version
go run ./cmd/pfm init

# Add a transaction
go run ./cmd/pfm add --type expense --amount 12.34 --category Food --note "coffee"

# List transactions
go run ./cmd/pfm list
go run ./cmd/pfm list --limit 5