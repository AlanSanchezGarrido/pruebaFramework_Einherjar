//repository la unica capa que sabe que SQL. recibe y regresa model.customer

package repository

import (
	"clientes-api/internal/customer/model"
	"context"

	postgres "code.nochebuena.dev/einherjar/db-postgres"
)
//contrato de persistencia. Lo consume service.customer, que solo conoce esta interfaz 
type Customer interface {
	Create(ctx context.Context,c model.Customer) error
	List(ctx context.Context)([]model.Customer,error)
	Get(ctx context.Context,id string)(model.Customer,error)
	Update(ctx context.Context,c model.Customer) error
	Delete(ctx context.Context, id string) error
	FindByEmail (ctx context.Context,email string)(model.Customer,error)
}
//estructra que implementa la interfaz contiene el campo de la interfaz de del paquete db-postgres
type customer struct{
	db postgres.Provider
}

//Falla en tiempo de compilacion si customer deja de cumplir el contrato(si falta un metodo)
var _ Customer = (*customer)(nil) 
//NewCustomer recibe postgres.Provider, entonces el repositori solo necesita ejecutar queries,
//no arrancar ni apagar la bd
func NewCustomer(db postgres.Provider) Customer {
	return &customer{db: db}
}


func (c *customer)Create(ctx context.Context, t model.Customer)error  {
	const q = `INSERT INTO customers (id,name,lastname,cellphone,email,address,password_hash,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
    // GetExecutor regresa la transacción activa si venimos dentro de un
	// UnitOfWork, o el pool si no. El repository no decide eso ni le importa.
	_,err := c.db.GetExecutor(ctx).Exec(ctx,q,t.ID,t.Name,t.LastName,t.CellPhone,t.Email,t.Address,t.PasswordHash,t.CreatedAt,t.UpdatedAt)
	if err!= nil {
		//HandlerError traduce codigo de postgres a xerros:UNIQUE → ErrAlreadyExists,
		//// sql.ErrNoRows → ErrNotFound, etc. httputil los mapea al status HTTP correcto.
		return c.db.HandleError(err)
	}
	return nil
}

func(c *customer)List(ctx context.Context)([]model.Customer,error) {
	const q = `SELECT id,name,lastname,cellphone,email,address,created_at,updated_at FROM customers ORDER BY created_at DESC`
	
	rows,err := c.db.GetExecutor(ctx).Query(ctx,q)
	
	if err!= nil {
		return nil,c.db.HandleError(err)
	}
    //Cerramos y liberamos los datos despues de que la funcion actual termine 
	defer rows.Close()
    // Slice inicializado, no nil: así la API responde [] y no null cuando no hay nada.
	out := []model.Customer{}

	for rows.Next(){
		t,err := scanCustomer(rows)

		if err!=nil {
			return nil,c.db.HandleError(err)
		}
		out = append(out, t)
	}
	// rows.Err() atrapa el error que rows.Next() se traga al regresar false.
	if err := rows.Err(); err != nil {
		return nil, c.db.HandleError(err)
	}
	return out, nil
}

func (c *customer)Get(ctx context.Context,id string) (model.Customer, error)  {
	const q = `SELECT id,name,lastname,cellphone,email,address,created_at,updated_at FROM customers WHERE id = $1`

	row,err := scanCustomer(c.db.GetExecutor(ctx).QueryRow(ctx,q,id))

	if err!=nil {
		// sql.ErrNoRows sale de aquí como ErrNotFound → 404, sin un if extra.
		return model.Customer{},c.db.HandleError(err)
	}

	return row,nil
}

func (c *customer)Update(ctx context.Context,t model.Customer) error  {
	const q=`UPDATE customers SET name=$1,lastname=$2,cellphone=$3,email=$4,address=$5,updated_at=$6 WHERE id=$7`

	_,err := c.db.GetExecutor(ctx).Exec(ctx,q,t.Name,t.LastName,t.CellPhone,t.Email,t.Address,t.UpdatedAt,t.ID)
	if err!=nil {
		return c.db.HandleError(err)
	}
	return nil
}

func (c *customer)Delete(ctx context.Context, id string)error  {
	const q =`DELETE FROM customers WHERE id=$1`

	res,err := c.db.GetExecutor(ctx).Exec(ctx,q,id)

	if err!= nil {
		return c.db.HandleError(err)
	}
	return requireOneRow(res,id)

}

//funcion para verificar el email con el cual se realizara el inicio de sesion 
func (c *customer) FindByEmail(ctx context.Context, email string) (model.Customer, error) {
	const q = `SELECT id, name, lastname, cellphone, email, address, password_hash, created_at, updated_at FROM customers WHERE email=$1`

	row, err := scanCustomerWithCustomer(c.db.GetExecutor(ctx).QueryRow(ctx, q, email))

	if err != nil {
		return model.Customer{}, c.db.HandleError(err)
	}

	return row, nil
}




