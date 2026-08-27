package service

import (
	"clientes-api/internal/customer/repository"
	"clientes-api/internal/perm"
	"clientes-api/internal/auth/util"
	"context"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/xerrors"
)



type Auth struct {
	logger logging.Logger
	repo   repository.Customer
	signer authjwt.Signer
	cfg    authjwt.TokenConfig
}

func NewAuth(logger logging.Logger, repo repository.Customer, signer authjwt.Signer, cfg authjwt.TokenConfig) *Auth {
	return &Auth{logger: logger, repo: repo, signer: signer, cfg: cfg}
}

func (a *Auth) Login(ctx context.Context, email, password string) (authjwt.TokenPair,error) {
	cust, err := a.repo.FindByEmail(ctx, email)

	if err != nil {
		a.logger.WithContext(ctx).Warn("login fallido", "email", email, "motivo", "password")
		return  authjwt.TokenPair{}, xerrors.Unauthorized("credenciales inválidas")
	}

	if err := util.CheckPassword(cust.PasswordHash, password); err != nil {
		a.logger.WithContext(ctx).Warn("login fallido", "email", email, "motivo", "password")
		return authjwt.TokenPair{}, xerrors.Unauthorized("credenciales inválidas")
	}

	mask := security.PermissionMask(0).
		Grant(perm.ReadCustomers).
		Grant(perm.WriteCustomers).
		Grant(perm.ReadShoppings).
		Grant(perm.WriteShoppings)

	claims := map[string]any{
		"perms": map[string]any{
			"customers": int64(mask),
			"shopping":  int64(mask),
		},
	}

	pair, err := authjwt.IssueTokenPair(a.signer, cust.ID, claims, a.cfg)
	
	if err != nil {
		return authjwt.TokenPair{}, err
	}

	a.logger.Info("login exitoso", "customer_id", cust.ID)
	return pair,nil
}
