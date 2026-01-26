package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gBogdan95/ATAD---Project-Personal-Finance-CLI-Manager/internal/db"
	"github.com/gBogdan95/ATAD---Project-Personal-Finance-CLI-Manager/internal/models"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("pfm v" + version)

	case "hello":
		fmt.Println("Hello! PFM is working.")

	case "init":
		if err := cmdInit(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

	case "add":
		if err := cmdAdd(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

	case "list":
		if err := cmdList(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("pfm - Personal Finance CLI Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  pfm version")
	fmt.Println("  pfm hello")
	fmt.Println("  pfm init")
	fmt.Println("  pfm add --type expense --amount 12.34 --category Food [--date YYYY-MM-DD] [--note \"...\"]")
	fmt.Println("  pfm list [--limit 20]")
}

func cmdInit() error {
	_, err := db.EnsureDir()
	if err != nil {
		return err
	}

	path, err := db.DBPath()
	if err != nil {
		return err
	}

	conn, err := db.Open(path)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := db.InitSchema(conn); err != nil {
		return err
	}

	fmt.Println("Initialized database at:", path)
	return nil
}

func cmdAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	typ := fs.String("type", "", "income or expense (required)")
	amount := fs.String("amount", "", "amount like 12.34 (required)")
	category := fs.String("category", "", "category name (required)")
	date := fs.String("date", "", "YYYY-MM-DD (default: today)")
	note := fs.String("note", "", "optional note")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *typ != "income" && *typ != "expense" {
		return fmt.Errorf("--type must be income or expense")
	}
	if strings.TrimSpace(*amount) == "" || strings.TrimSpace(*category) == "" {
		return fmt.Errorf("--amount and --category are required")
	}

	d := strings.TrimSpace(*date)
	if d == "" {
		d = time.Now().Format("2006-01-02")
	} else {
		if _, err := time.Parse("2006-01-02", d); err != nil {
			return fmt.Errorf("invalid --date (use YYYY-MM-DD): %w", err)
		}
	}

	cents, err := parseAmountToCents(*amount)
	if err != nil {
		return err
	}

	_, err = db.EnsureDir()
	if err != nil {
		return err
	}
	path, err := db.DBPath()
	if err != nil {
		return err
	}

	conn, err := db.Open(path)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := db.InitSchema(conn); err != nil {
		return err
	}

	id, err := db.InsertTransaction(conn, models.Transaction{
		Date:        d,
		Type:        *typ,
		AmountCents: cents,
		Category:    *category,
		Note:        *note,
	})
	if err != nil {
		return err
	}

	fmt.Println("Added transaction with ID:", id)
	return nil
}

func parseAmountToCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("amount is empty")
	}

	parts := strings.SplitN(s, ".", 2)

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %w", err)
	}

	var frac int64
	if len(parts) == 2 {
		f := parts[1]
		if len(f) == 1 {
			f += "0"
		}
		if len(f) > 2 {
			f = f[:2]
		}
		frac, err = strconv.ParseInt(f, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount: %w", err)
		}
	}

	return whole*100 + frac, nil
}

func cmdList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	limit := fs.Int("limit", 20, "number of transactions to show")
	if err := fs.Parse(args); err != nil {
		return err
	}

	_, err := db.EnsureDir()
	if err != nil {
		return err
	}
	path, err := db.DBPath()
	if err != nil {
		return err
	}

	conn, err := db.Open(path)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := db.InitSchema(conn); err != nil {
		return err
	}

	items, err := db.ListTransactions(conn, *limit)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println("No transactions found.")
		return nil
	}

	fmt.Printf("Showing last %d transactions:\n", len(items))
	for _, t := range items {
		fmt.Printf("#%d  %s  %-7s  %s  %s",
			t.ID,
			t.Date,
			t.Type,
			formatCents(t.AmountCents),
			t.Category,
		)
		if strings.TrimSpace(t.Note) != "" {
			fmt.Printf("  (%s)", t.Note)
		}
		fmt.Println()
	}

	return nil
}

func formatCents(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d.%02d", sign, cents/100, cents%100)
}
