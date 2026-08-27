package handler

import (
	"clientes-api/internal/Compras/dto"
	"clientes-api/internal/Compras/service"
	"context"
	"net/http"

	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/valid"
	"code.nochebuena.dev/einherjar/web/httputil"
	"github.com/go-chi/chi/v5"
)

type Shopping interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type shopping struct {
	logger    logging.Logger
	validator valid.Validator
	srv       service.Shopping
}

var _ Shopping = (*shopping)(nil)

func NewShopping(logger logging.Logger, validator valid.Validator, srv service.Shopping) Shopping {
	return &shopping{
		logger:    logger,
		validator: validator,
		srv:       srv,
	}
}

func (c *shopping) Create(w http.ResponseWriter, r *http.Request) {
	httputil.Handle(c.validator, c.logger, func(ctx context.Context, in dto.RegisterShopping) (dto.Shopping, error) {
		t, err := c.srv.Create(ctx, in.Total, in.Article)
		if err != nil {
			return dto.Shopping{}, err
		}
		return dto.FromShopping(t), nil
	})(w, r)
}
func (c *shopping) List(w http.ResponseWriter, r *http.Request) {
	httputil.HandleNoBody(c.logger, func(ctx context.Context) (dto.ShoppingList, error) {
		ss, err := c.srv.ListMine(ctx)
		if err != nil {
			return dto.ShoppingList{}, err
		}
		return dto.FromShoppings(ss), nil
	})(w, r)
}

func (c *shopping) Get(w http.ResponseWriter, r *http.Request) {
	s, err := c.srv.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		httputil.Error(c.logger, w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FromShopping(s))
}

func (c *shopping) Update(w http.ResponseWriter, r *http.Request) {
	httputil.Handle(c.validator, c.logger, func(ctx context.Context, in dto.UpdatedShopping) (dto.Shopping, error) {
		s, err := c.srv.Update(ctx, chi.URLParam(r, "id"), in.Total, in.Article)

		if err != nil {
			return dto.Shopping{}, err
		}
		return dto.FromShopping(s), nil
	})(w, r)
}

func (c *shopping) Delete(w http.ResponseWriter, r *http.Request) {
	if err := c.srv.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		httputil.Error(c.logger, w, r, err)
		return
	}
	httputil.NoContent(w)
}
