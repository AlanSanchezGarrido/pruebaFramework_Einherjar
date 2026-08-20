// package deto son las formas que viajan por HTTP: lo que entra en el body y lo que sale en ela respuesta es el contrato del cliente
//
// Esta separo de model.Customer a proposito si en dado caso en el futuro se le agrega una columna a la tbala, model cambia y el contrati
// de la api no-mientras nadie toque este paquete, el cliente no se entera
//
// aqui tambien viven los mappers, para que handler quede de puro transporte
// dto conoce a model; model no conoce a dto
package dto

import (
	"clientes-api/internal/model"
	"time"
)

//createCustomer es el body del POST/api/v1/customer
//
//los tags `validate` los ejecuta core/valid de httputil.Handler,antes de que el handler llame al service : si algo falla, el service nunca
// se entera y el cliente recibe un 400 con el campo que reprobo

type CreateCustomer struct {
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	CellPhone string `json:"cellphone" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Address   string `json:"address" validate:"required"`
}

//Tomodel arma la entidad con la que el cliente puede recibir. El ID y las fechas no pertenecen aqui si no al service, por eso no se responden

func (r CreateCustomer) Tomodel() model.Customer {
	return model.Customer{Name: r.Name, LastName: r.LastName, CellPhone: r.CellPhone, Email: r.Email, Address: r.Address}
}

type UpdatedCustomer struct {
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	CellPhone string `json:"cellphone" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Address   string `json:"address" validate:"required"`
}

func (r UpdatedCustomer) Tomodel(id string) model.Customer {
	return model.Customer{Name: r.Name, LastName: r.LastName, CellPhone: r.CellPhone, Email: r.Email, Address: r.Address}
}

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastname"`
	CellPhone string `json:"cellphone"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

func FromCustomer(t model.Customer) Customer  {
	return Customer{
		ID: t.ID,
		Name: t.Name,
		LastName: t.LastName,
		CellPhone: t.CellPhone,
		Email: t.Email,
		Address: t.Address,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt}
}

type CustomerList struct{
	Customer []Customer `json:"customer"`
	Count int `json:"count"` 
}

func FromCustomers(cs []model.Customer) CustomerList  {
	customers := make([]Customer,0,len(cs))
	for _, v := range cs {
		customers = append(customers,FromCustomer(v))
	}
	return CustomerList{Customer: customers,Count: len(customers)}
}




