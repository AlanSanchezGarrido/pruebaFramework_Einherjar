package handler

import (
	"clientes-api/internal/auth/dto"
	"clientes-api/internal/auth/service"
	"context"
	"net/http"

	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/valid"
	"code.nochebuena.dev/einherjar/web/httputil"
)


type Auth struct{
	logger logging.Logger
	validator valid.Validator
	svc *service.Auth
}

func NewAuth(logger logging.Logger,validator valid.Validator,svc *service.Auth) *Auth {
	return &Auth{logger: logger, validator: validator, svc: svc}
}

func (a *Auth)Login(w http.ResponseWriter, r *http.Request)  {
	httputil.Handle(a.validator,a.logger,func(ctx context.Context, in dto.Login) (dto.AuthResponse,error) {
		pair,err:=a.svc.Login(ctx,in.Email,in.Password)
		if err!=nil {
			return dto.AuthResponse{},err
		}
		return dto.FromAuth(pair.AccessToken,pair.RefreshToken,pair.ExpiresIn),nil

	})(w,r)
}
	