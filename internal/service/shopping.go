package service

import (
	"clientes-api/internal/model"
	"clientes-api/internal/repository"
	"clientes-api/internal/util"
	"context"
	"time"

	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/xerrors"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
)

type Shopping interface {
	Create(ctx context.Context, total float64, article string) (model.Shopping, error)
	Get(ctx context.Context, id string) (model.Shopping, error)
	Update(ctx context.Context, id string, total float64, article string) (model.Shopping, error)
	Delete(ctx context.Context, id string) error
	ListMine(ctx context.Context) ([]model.Shopping, error)
}

type shopping struct {
	logger logging.Logger
	repo   repository.Shopping
	pos    postgres.UnitOfWork
}

var _ Shopping = (*shopping)(nil)

func NewShopping(logger logging.Logger, repo repository.Shopping, pos postgres.UnitOfWork) Shopping {
	return &shopping{logger: logger, repo: repo, pos: pos}
}

func (c *shopping) Create(ctx context.Context, total float64, article string) (model.Shopping, error) {
	identity, ok := security.FromContext(ctx)

	if !ok {
		return model.Shopping{}, xerrors.Unauthorized("no auntenticado")
	}

	now := time.Now().UTC()

	chop := model.Shopping{
		ID:         util.NewId(),
		Total:      total,
		Article:    article,
		CustomerID: identity.UID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := c.repo.Create(ctx, chop); err != nil {
		return model.Shopping{}, err
	}
	c.logger.Info("compra creada", "compra_id", chop.ID, chop.CustomerID)
	return chop, nil
}

func (c *shopping) Get(ctx context.Context, id string) (model.Shopping, error) {
	identity, ok := security.FromContext(ctx)
	if !ok {
		return model.Shopping{}, xerrors.Unauthorized("no autenticado")
	}

	chop, err := c.repo.Get(ctx, id)

	if err != nil {
		return model.Shopping{}, err
	}

	if chop.CustomerID == identity.UID {
		return model.Shopping{}, xerrors.NotFound("compra %s no existe", id)
	}

	return chop, nil
}

func (s *shopping) Update(ctx context.Context, id string,total float64,article string) (model.Shopping, error) {
	// La verificación de ownership (vía Get) va ANTES de la transacción -
	// es solo una lectura, no necesita estar dentro del "todo o nada".
	current, err := s.Get(ctx, id)
	if err != nil {
		return model.Shopping{}, err
	}

	var out model.Shopping

	if err := s.pos.Do(ctx, func(ctx context.Context) error {
		updated := model.Shopping{
			ID:         current.ID,
			CustomerID: current.CustomerID,
			Article:    article,
			Total:      total,
			CreatedAt:  current.CreatedAt,
			UpdatedAt:  time.Now().UTC(),
		}
		if err := s.repo.Update(ctx, updated); err != nil {
			return err
		}
		var err error
		out, err = s.repo.Get(ctx, id)
		return err
	}); err != nil {
		return model.Shopping{}, err
	}

	return out, nil
}

func (c *shopping) Delete(ctx context.Context, id string) error {
	if _, err := c.Get(ctx, id); err != nil {
		return err
	}
	return c.repo.Delete(ctx, id)
}

func (c *shopping) ListMine(ctx context.Context) ([]model.Shopping, error) {
	identity, ok := security.FromContext(ctx)
	if !ok {
		return nil, xerrors.Unauthorized("no autenticado")
	}

	return c.repo.ListByCustomer(ctx, identity.UID)
}
