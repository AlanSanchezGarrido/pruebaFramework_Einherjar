package handler

import (
	"clientes-api/internal/service"
	"clientes-api/internal/dto"
	"context"
	"net/http"


	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/valid"
	"code.nochebuena.dev/einherjar/web/httputil"
)

type Auth interface {
	Login(w http.ResponseWriter, r *http.Request)
}
type auth struct{
	logger logging.Logger
	validator valid.Validator
	svc service.Auth
}

var _ Auth = (*auth)(nil)
func NewAuth(logger logging.Logger,validator valid.Validator,svc service.Auth) Auth {
	return &auth{logger: logger, validator: validator, svc: svc}
}

func (a *auth)Login(w http.ResponseWriter, r *http.Request)  {
	httputil.Handle(a.validator,a.logger,func(ctx context.Context, in dto.Login) (dto.AuthResponse,error) {
		cust,pair,err:=a.svc.Login(ctx,in.Email,in.Password)
		if err!=nil {
			return dto.AuthResponse{},err
		}
		return dto.FromAuth(cust,pair.AccessToken,pair.RefreshToken,pair.ExpiresIn),nil

	})(w,r)
}