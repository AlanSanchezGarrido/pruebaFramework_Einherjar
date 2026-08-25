package service

import (
	"clientes-api/internal/model"
	"clientes-api/internal/repository"
	"clientes-api/internal/util"
	"clientes-api/internal/perm"
	"context"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/xerrors"
)

type Auth interface {
	Login(ctx context.Context, email, password string)(model.Customer,authjwt.TokenPair,error)
}

type auth struct{
	logger logging.Logger
	repo repository.Customer
	signer authjwt.Signer
	cfg authjwt.TokenConfig
}

var _ Auth = (*auth)(nil)

func NewAuth(logger logging.Logger,repo repository.Customer,signer authjwt.Signer, cfg authjwt.TokenConfig) Auth {
	return &auth{logger: logger, repo: repo, signer: signer, cfg: cfg}
}

func (a *auth)Login (ctx context.Context, email, password string) (model.Customer,authjwt.TokenPair,error)  {
	cust,err:=a.repo.FindByEmail(ctx,email)

	if err!=nil {
		a.logger.WithContext(ctx).Warn("login fallido","email",email,"motivo","password")
		return model.Customer{}, authjwt.TokenPair{}, xerrors.Unauthorized("credenciales inválidas")
	}

	if err:= util.CheckPassword(cust.PasswordHash,password);err!=nil {
		a.logger.WithContext(ctx).Warn("login fallido","email",email,"motivo","password")
	return model.Customer{}, authjwt.TokenPair{}, xerrors.Unauthorized("credenciales inválidas")
	}
	
	mask :=security.PermissionMask(0).Grant(perm.ReadCustomers).Grant(perm.WriteCustomers)
	claims := map[string]any{
		"perms": map[string]any{
			"customers": int64(mask),
		},
	}

	pair,err:=authjwt.IssueTokenPair(a.signer, cust.ID,claims,a.cfg)
	if err!=nil {
		return model.Customer{},authjwt.TokenPair{},err
	}
	a.logger.Info("login exitoso","customer_id",cust.ID)
	return cust,pair,nil
}