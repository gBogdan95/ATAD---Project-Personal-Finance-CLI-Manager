package models

type Transaction struct {
	ID          int64
	Date        string
	Type        string
	AmountCents int64
	Category    string
	Note        string
}
