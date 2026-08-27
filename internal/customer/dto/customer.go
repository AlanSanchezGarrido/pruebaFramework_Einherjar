// package deto son las formas que viajan por HTTP: lo que entra en el body y lo que sale en ela respuesta es el contrato del cliente
//
// Esta separo de model.Customer a proposito si en dado caso en el futuro se le agrega una columna a la tbala, model cambia y el contrati
// de la api no-mientras nadie toque este paquete, el cliente no se entera
//
// aqui tambien viven los mappers, para que handler quede de puro transporte
// dto conoce a model; model no conoce a dto
package dto

import (
	"clientes-api/internal/customer/model"
	"time"
)

// createCustomer es el body del POST/api/v1/customer
//
// los tags `validate` los ejecuta core/valid de httputil.Handler,antes de que el handler llame al service : si algo falla, el service nunca
// se entera y el cliente recibe un 400 con el campo que reprobo
// molde exacto que revisara y contendra los datos del cliente que quiere registrarse en el sistema
type CreateCustomer struct {
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	CellPhone string `json:"cellphone" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Address   string `json:"address" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

// Tomodel arma la entidad con la que el cliente puede recibir. El ID y las fechas no pertenecen aqui si no al service, por eso no se responden
// crea un formato que utiliza el negocio,extrae solo la informacion valida del formulario y crea la entidad de negocio oficial
func (r CreateCustomer) Tomodel() model.Customer {
	return model.Customer{Name: r.Name, LastName: r.LastName, CellPhone: r.CellPhone, Email: r.Email, Address: r.Address}
}

type UpdatedCustomer struct {
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	CellPhone string `json:"cellphone" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Address   string `json:"address" validate:"required"`
}

func (r UpdatedCustomer) Tomodel(id string) model.Customer {
	return model.Customer{ID: id, Name: r.Name, LastName: r.LastName, CellPhone: r.CellPhone, Email: r.Email, Address: r.Address}
}


//estructura de respuesta Json al crean un cliente nuevo
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	CellPhone string    `json:"cellphone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}
//funcion para realizar la estructura y regresar un cliente
func FromCustomer(t model.Customer) Customer {
	return Customer{
		ID:        t.ID,
		Name:      t.Name,
		LastName:  t.LastName,
		CellPhone: t.CellPhone,
		Email:     t.Email,
		Address:   t.Address,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt}
}

//Estructura para respues de una lista de clientes.regresa los clientes y el numero 
type CustomerList struct {
	Customer []Customer `json:"customer"`
	Count    int        `json:"count"`
}

// FromTodos usa make con len 0 y no var: un slice nil se serializa como null, y
// la llave todos tiene que traer [] cuando no hay nada, no null.
func FromCustomers(cs []model.Customer) CustomerList {
	customers := make([]Customer, 0, len(cs))
	for _, v := range cs {
		customers = append(customers, FromCustomer(v))
	}
	return CustomerList{Customer: customers, Count: len(customers)}
}
