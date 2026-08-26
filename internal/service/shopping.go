package service

import (
	"clientes-api/internal/model"
	"clientes-api/internal/repository"
	"context"

	"code.nochebuena.dev/einherjar/contracts/logging"
)

type Shopping interface {
	Create(ctx context.Context,total float64)(model.Shopping,error)
	Get(ctx context.Context, id string)(model.Shopping,error)
	Update(ctx context.Context,id string,total float64)(model.Shopping,error)
	Delete(ctx context.Context,id string)error
	ListMine(ctx context.Context)([]model.Shopping,error)
}

type shopping struct{
	logger logging.Logger
	repo repository.Shopping

}