//molde de datos, no conoce sql, ni http
//no contiene tags a proposito 
package model

import "time"
//customer es la entidad que viaja desde service a repositori y de regreso dto la traduce 
type Customer struct {
	ID        string
	Name      string
	LastName  string
	CellPhone string
	Email     string
	Address    string
	PasswordHash string
	CreatedAt time.Time
	UpdatedAt time.Time
}

