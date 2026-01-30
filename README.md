# Personal Finance CLI Manager (pfm)

A command-line app for tracking income and expenses locally using SQLite.

This project is developed for the **Advanced Technologies for Application Development** course (Master’s).
It implements the “Personal Finance CLI Manager” project: CLI subcommands, local storage, reporting, and budgeting.

## Current Features (Architecture + Basic Functionality)
- Local SQLite database storage (`init`)
- Transactions:
  - `add` inserts income/expense transactions
  - `list` prints recent transactions (default 20, configurable via `--limit`)
- Reports:
  - `report --month YYYY-MM` prints monthly income/expense/net totals
- Budgets:
  - `budget set` creates/updates a monthly budget for a category
  - `budget status` shows spending vs budget (OK/OVER)

## Run

### Quick start
```bash
go run ./cmd/pfm hello  
go run ./cmd/pfm version  
go run ./cmd/pfm init  
```

## Add transactions
# Expense
```bash
go run ./cmd/pfm add --type expense --amount 12.34 --category Food --note "coffee"  
```

# Income  
```bash
go run ./cmd/pfm add --type income --amount 1000 --category Salary --note "test salary"  
```

## List transactions
```bash
go run ./cmd/pfm list  
go run ./cmd/pfm list --limit 5  
```

## Monthly report
```bash
go run ./cmd/pfm report --month 2026-01  
```

## Budgets
# Set/update budget for a category in a month  
```bash
go run ./cmd/pfm budget set --month 2026-01 --category Food --amount 50  
```

# Show budget status for the month  
```bash
go run ./cmd/pfm budget status --month 2026-01  
```

## Architecture
- `cmd/pfm/` - CLI entrypoint
- `internal/db/` - SQLite connection, schema, and queries
- `internal/models/` - domain models (Transaction, Budget, ...)

## Database
The SQLite database is stored in the OS user config directory under a `pfm` folder.  
Example (Windows): `C:\Users\<you>\AppData\Roaming\pfm\pfm.db`

## Notes / Assumptions
- Dates are stored as text in `YYYY-MM-DD` format.
- Amounts are stored as integer cents to avoid floating-point errors.

## Requirements roadmap
- [x] CLI app with subcommands
- [x] Local SQLite persistence
- [x] Budgets with alerts/status
- [x] Reports (monthly totals)