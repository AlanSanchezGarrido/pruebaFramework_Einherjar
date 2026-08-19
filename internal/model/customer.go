package model

import "time"

type Customer struct {
	ID        string
	Name      string
	LastName  string
	CellPhone string
	Email     string
	Address    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

