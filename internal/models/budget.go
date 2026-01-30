package models

type Budget struct {
	ID         int64
	Month      string
	Category   string
	LimitCents int64
}
