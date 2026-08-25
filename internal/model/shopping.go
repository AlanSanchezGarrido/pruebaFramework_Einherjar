package model

import "time"

type Shopping struct {
	ID          string
	Total       string
	Article     string
	Customer_id string
	Create_At   time.Time
	Update_At   time.Time
}
