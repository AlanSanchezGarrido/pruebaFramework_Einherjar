package dto

import (
	"clientes-api/internal/Compras/model"
	"time"
)

// structura de json para crear una compra
type RegisterShopping struct {
	Total   float64 `json:"total" validate:"required"`
	Article string `json:"article" validate:"required"`
}

// estructura json para actualizar una compra
type UpdatedShopping struct {
	Total   float64 `json:"total" validate:"required"`
	Article string `json:"article" validate:"required"`
}

// estructura json para responder una compra
type Shopping struct {
	ID         string    `json:"id"`
	Article    string    `json:"article"`
	Total      float64   `json:"total"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// funcion para regresar esa estructura json
func FromShopping(s model.Shopping) Shopping {
	return Shopping{
		ID:         s.ID,
		Article:    s.Article,
		Total:      s.Total,
		CustomerID: s.CustomerID,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

//estructura para responder una lista de comrpas
type ShoppingList struct {
	Shopping []Shopping `json:"shopping"`
	Count    int        `json:"count"`
}

//funcion para regresar una lista de compras con contador
func FromShoppings(ss []model.Shopping) ShoppingList  {
	out:=make([]Shopping,0,len(ss))

	for _, v := range ss {
		out = append(out, FromShopping(v))
	}
	return ShoppingList{Shopping: out,Count: len(out)}
}
