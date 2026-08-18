package repository

import (
	"clientes-api/internal/model"

	"code.nochebuena.dev/einherjar/core/xerrors"
	"github.com/jackc/pgx/v5/pgconn"
)

type scanner interface {
	Scan(dest ...any) error
}

func scanCustomer(s scanner) (model.Customer, error) {
var(
	c 	model.Customer
)

if err := s.Scan(&c.ID,&c.Name,&c.LastName,&c.CellPhone,&c.CellPhone,&c.Email,&c.Address,&c.CreatedAt,&c.UpdatedAt);err!=nil {
	return model.Customer{},err
}
return c,nil
}

func requireOneRow(tag pgconn.CommandTag, id string) error {
	if tag.RowsAffected() ==0 {
		return xerrors.NotFound("cliente %s no existe",id)
	}
	return nil
}

