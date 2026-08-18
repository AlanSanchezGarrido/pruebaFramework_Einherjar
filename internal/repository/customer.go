package repository

import (
	"clientes-api/internal/model"
	"context"

	postgres "code.nochebuena.dev/einherjar/db-postgres"
)

type Customer interface {
	Create(ctx context.Context,c model.Customer) error
	List(ctx context.Context)([]model.Customer,error)
	Get(ctx context.Context,id string)(model.Customer,error)
	Updated(ctx context.Context,c model.Customer) error
	Deleted(ctx context.Context, id string) error

}

type customer struct{
	db postgres.Provider
}

var _ Customer = (*customer)(nil) 

func NewCustomer(db postgres.Provider) Customer {
	return &customer{db: db}
}

func (c *customer)Create(ctx context.Context, t model.Customer)error  {
	const q = `INSERT INTO customers (id,name,lastname,cellphone,email,address,created_at,updated_at) VALUE (?,?,?,?,?,?,?,?)`

	_,err := c.db.GetExecutor(ctx).Exec(ctx,q,t.ID,t.Name,t.LastName,t.CellPhone,t.Email,t.Address,t.CreatedAt,t.UpdatedAt)
	if err!= nil {
		return c.db.HandleError(err)
	}
	return nil
}

func(c *customer)List(ctx context.Context)([]model.Customer,error) {
	const q = `SELECT (id,name,lastname,cellphone,email,address,created_at,updated_at) FROM customers ORDER BY DESC`
	
	rows,err := c.db.GetExecutor(ctx).Query(ctx,q)
	
	if err!= nil {
		return nil,c.db.HandleError(err)
	}

	defer rows.Close()

	out := []model.Customer{}

	for rows.Next(){
		t,err := scanCustomer(rows)

		if err!=nil {
			return nil,c.db.HandleError(err)
		}
		out = append(out, t)
	}
	return out,nil
}

func (c *customer)Get(ctx context.Context,id string) (model.Customer, error)  {
	const q = `SELECT id,name,lastname,cellphone,email,address,created_at,updated_at FROM customers WHERE id = ?`

	row,err := scanCustomer(c.db.GetExecutor(ctx).QueryRow(ctx,q,id))

	if err!=nil {
		return model.Customer{},c.db.HandleError(err)
	}

	return row,nil
}

func (c *customer)Updated(ctx context.Context,t model.Customer) error  {
	const q=`UPDATE customers SET name=?,lastname=?,cellphone=?,email=?,address=?,update_at=? WHERE id=?`

	_,err := c.db.GetExecutor(ctx).Exec(ctx,q,t.Name,t.LastName,t.CellPhone,t.Email,t.Address,t.UpdatedAt)
	if err!=nil {
		return c.db.HandleError(err)
	}
	return nil
}

func (c *customer)Deleted(ctx context.Context, id string)error  {
	const q =`DELETE FROM customers WHERE id=?`

	res,err := c.db.GetExecutor(ctx).Exec(ctx,q,id)

	if err!= nil {
		return c.db.HandleError(err)
	}

	return requireOneRow(res,id)

}



