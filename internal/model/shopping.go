package model

import "time"

type Shopping struct {
	ID         string
	Total      float64
	Article    string
	CustomerID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}