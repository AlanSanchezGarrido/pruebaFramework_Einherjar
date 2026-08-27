package repository

import (
	"clientes-api/internal/Compras/model"

	"code.nochebuena.dev/einherjar/core/xerrors"
	"github.com/jackc/pgx/v5/pgconn"
)

type scan interface {
	Scan(dest ...any) error
}

func scanCompra(s scan) (model.Shopping, error) {
	var c model.Shopping

	if err := s.Scan(&c.ID, &c.Total, &c.Article, &c.CustomerID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return model.Shopping{}, nil
	}
	return c, nil
}

func requireOneRow(tag pgconn.CommandTag,id string)error  {
	if tag.RowsAffected()==0 {
		return xerrors.NotFound("no existe el usuario con ese %s",id)
	}
	return nil
}