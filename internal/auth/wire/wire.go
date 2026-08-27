package wire

import (
	"clientes-api/internal/auth/handler"
	"clientes-api/internal/auth/service"
	"clientes-api/internal/customer/repository"
	"clientes-api/internal/identity"
	"net/http"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"
	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/auth/rbac"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
	"github.com/go-chi/chi/v5"
)

// Wire construye todo lo relacionado a auth y devuelve el PermissionProvider
// para que otros módulos (customer, shopping) lo usen en su AuthzMiddleware.
func Wire(lc launcher.Launcher,srv server.Server, logger logging.Logger, validator valid.Validator, db postgres.Component, jwtCfg authjwt.TokenConfig, jwtSecret string) (func(http.Handler) http.Handler, authmw.IdentityEnricher, security.PermissionProvider) {
	customerRepo := repository.NewCustomer(db)

	signer := authjwt.NewHMACSigner([]byte(jwtSecret))
	enricher := identity.NewEnricher(customerRepo)
	permissions := rbac.NewClaimsPermissionProvider("perms", authmw.GetClaims)

	// Globales — deben ir ANTES de que otros módulos registren sus rutas.
	authMiddleware := authjwt.AuthMiddleware(logger, signer, nil)

	authSRV := service.NewAuth(logger,customerRepo, signer, jwtCfg)
	authH := handler.NewAuth(logger, validator, authSRV)

	lc.BeforeStart(func() error {
		srv.Route("/api/v1/auth", func(r chi.Router) {
			r.Post("/login", authH.Login)
		})
		return nil
	})

	return authMiddleware,enricher,permissions
}