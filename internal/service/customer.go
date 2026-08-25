// services vive toda la ogica del negocio, no conoce que es sql, no conoce las bases de datos., ni conoce lo que es
// HTTP status code , entonces solo habla con model osea el modelo
package service

import (
	"clientes-api/internal/model"
	"clientes-api/internal/repository"
	"clientes-api/internal/util"
	"context"
	"strings"
	"time"

	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/xerrors"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
)

//todo el contrato de negocio que consume handler.customer
//recibe y regresa model.customer y nada de customer.dto

type Customer interface{
	Create(ctx context.Context, t model.Customer, plainPassword string)(model.Customer,error)
	List(ctx context.Context)([]model.Customer,error)
	Get(ctx context.Context, id string)(model.Customer, error)
	Update(ctx context.Context,t model.Customer)(model.Customer,error)
	Delete(ctx context.Context,id string)error
}

type customer struct{
	logger logging.Logger
	repo repository.Customer
	pos postgres.UnitOfWork
}

var _ Customer= (*customer)(nil)

func NewCustomer(logger logging.Logger,repo repository.Customer,pos postgres.UnitOfWork) Customer {
	return &customer{logger: logger, repo: repo, pos: pos}
}

func (c *customer) Create(ctx context.Context, t model.Customer, plainPassword string)(model.Customer, error)  {
	 cus,err:=CleanInput(t)

	 if err!=nil {
		return model.Customer{},err
	 }

	 hash,err:=util.HasPassword(plainPassword)

	if err!=nil {
		return model.Customer{}, err
	}

	 now :=time.Now().UTC()
	 custom := model.Customer{
		ID: util.NewId(),
		Name: cus.Name,
		LastName: cus.LastName,
		CellPhone: cus.CellPhone,
		Email: cus.Email,
		Address: cus.Address,
		PasswordHash: hash,
		CreatedAt: now,
		UpdatedAt: now,
	 }
	
	if err := c.repo.Create(ctx,custom);err!=nil {
		return model.Customer{},err
	}

	c.logger.Info("Customer creado","customer_id",custom.ID)
	return custom,nil
}

func (c *customer) List(ctx context.Context)([]model.Customer,error)  {
	return	c.repo.List(ctx)
}



func (c *customer) Get(ctx context.Context,id string) (model.Customer, error)  {
	return c.repo.Get(ctx,id)
}

// Update es el único caso que toca la base dos veces: actualiza y relee para
// regresar el registro fresco. Las dos van dentro del mismo uow.Do, o sea la
// misma transacción — si el Get truena, el UPDATE se revierte.
//
// Fíjense que adentro del closure se usa el ctx del parámetro (el que sombrea al
// de afuera): ese trae la transacción inyectada. Usar el de afuera te saca de la
// transacción sin que nada truene, y ese bug no se ve.
func (c *customer) Update(ctx context.Context,t model.Customer)(model.Customer, error)  {
	cus,err:= CleanInput(t)

	if err!=nil {
		return model.Customer{},err
	}

	var out model.Customer

	if err := c.pos.Do(ctx,func(ctx context.Context) error {
		updated := model.Customer{
			ID: t.ID,
			Name: cus.Name,
			LastName: cus.LastName,
			CellPhone: cus.CellPhone,
			Email: cus.Email,
			Address: cus.Address,
			UpdatedAt: time.Now().UTC(),
		} 
		if err :=c.repo.Update(ctx,updated);err!=nil {
			return err
		}
		var err error
		out, err = c.repo.Get(ctx,t.ID)
		return err
	}); err!=nil {
		return	model.Customer{},err
	}
	c.logger.Info("cliente actualizado","customer_id",t.ID)
	return out,nil
}


func (c *customer) Delete(ctx context.Context,id string)error  {
	if err:=c.repo.Delete(ctx,id);err!=nil {
		return err
	}	
	c.logger.WithContext(ctx).Info("cliente eliminado","cliente_id",id)
	return nil
}



func CleanInput(t model.Customer)(model.Customer,error)  {
	name:= strings.TrimSpace(t.Name)
	lastname:= strings.TrimSpace(t.LastName)
	cellphone:= strings.TrimSpace(t.CellPhone)
	email:= strings.TrimSpace(t.Email)
	address:= strings.TrimSpace(t.Address)
	
	switch {
	case name=="":
		return model.Customer{},xerrors.InvalidInput("verificar el campo name")
	case lastname =="":
		return model.Customer{},xerrors.InvalidInput("verificar el campo lastname") 
	case cellphone == "":
		return model.Customer{},xerrors.InvalidInput("verificar el campo cellphone")
	case email =="":
		return model.Customer{},xerrors.InvalidInput("verificar el campo email")
	case address =="":
		return model.Customer{},xerrors.InvalidInput("verificar el campo address")
	}


	return model.Customer{ID: t.ID,Name: name,LastName: lastname,CellPhone: cellphone,Email: email,Address: address},nil
}