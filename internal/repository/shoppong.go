package repository

import (
	"clientes-api/internal/model"
	"context"

	postgres "code.nochebuena.dev/einherjar/db-postgres"
)

type Shopping interface {
	Create(ctx context.Context,s model.Shopping)error
	List(ctx context.Context)([]model.Shopping,error)
	Get(ctx context.Context,id string)(model.Shopping,error)
	Update(ctx context.Context,s model.Shopping)error
	Delete(ctx context.Context,id string)error
	ListByCustomer (ctx context.Context,customerid string)([]model.Shopping,error)
}

type shopping struct{
	pos postgres.Provider
}

var _ Shopping=(*shopping)(nil)

func (*shopping)NewShopping(pos postgres.Provider) Shopping {
	return &shopping{pos: pos}
}
//Funcion para crear una compra enlasada con un cliente
func (c *shopping)Create(ctx context.Context,s model.Shopping)error  {
	const q = `INSERT INTO (id,total,article,customer_id,create_at,update_at) VALUES ($1,$2,$3,$4,$5,$6)`

	_,err:=c.pos.GetExecutor(ctx).Exec(ctx,q,s.ID,s.Total,s.Article,s.Customer_id,s.Create_At,s.Update_At)

	if err!=nil {
		return c.pos.HandleError(err)
	}
	return nil
}


func (c *shopping)List(ctx context.Context)([]model.Shopping,error)  {
	const q = `SELECT id,total,article,customer_id,create_at,update_at From shopping`
	
	rows,err:=c.pos.GetExecutor(ctx).Query(ctx,q)
	if err!= nil {
		return nil,c.pos.HandleError(err)
	}

	defer rows.Close()

	shoppings:=[]model.Shopping{}

	for rows.Next(){
		sp,err:=scanCompra(rows)
		if err!=nil {
			return  nil,c.pos.HandleError(err)
		}
		shoppings = append(shoppings, sp)
	}
	if err:= rows.Err();err!=nil {
		return nil,c.pos.HandleError(err)
	}
	return shoppings,nil
}

func (c *shopping)Get(ctx context.Context,id string)(model.Shopping,error){
	const q =`SELECT id,total,article,customer_id,create_at,update_at From shopping WHERE id=$1`

	t,err:=scanCompra(c.pos.GetExecutor(ctx).QueryRow(ctx,q,id))
	if err!=nil {
		return model.Shopping{},c.pos.HandleError(err)
	}
	return t,nil
}

func (c *shopping)Update(ctx context.Context,s model.Shopping)error {
	const q=`UPDATE shopping SET total=$1,article=$2,update_at=$3 WHERE id=$4`

	res,err:=c.pos.GetExecutor(ctx).Exec(ctx,q,s.Total,s.Article,s.Update_At,s.ID)

	if err!=nil {
		return c.pos.HandleError(err)
	}
	return requireOneRow(res,s.ID)
}

func (c *shopping)Delete(ctx context.Context,id string)error  {
	const q = `DELETE FROM shopping WHERE id=$1`

	t,err:=c.pos.GetExecutor(ctx).Exec(ctx,q,id)

	if err!=nil {
		return c.pos.HandleError(err)
	}
	return  requireOneRow(t,id)
}

func (c *shopping)ListByCustomer(ctx context.Context,cusomerid string)([]model.Shopping,error)  {
	const q =  `SELECT id,customer_id,total,created_at,updated_at FROM compras WHERE customer_id=$1 ORDER BY created_at DESC`
	rows,err:=c.pos.GetExecutor(ctx).Query(ctx,q,cusomerid)

	if err!=nil {
		return nil,c.pos.HandleError(err)
	}
	defer rows.Close()

	chopings :=[]model.Shopping{}
	for rows.Next(){
		t,err := scanCompra(rows)
		if err!=nil {
			return nil,c.pos.HandleError(err)
		}
		chopings = append(chopings,t)
	} 
	if err:=rows.Err();err!=nil {
		return nil,c.pos.HandleError(err)
	}
	return chopings,nil
}