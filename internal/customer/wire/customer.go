package wire

import (
	"clientes-api/internal/customer/handler"
	"clientes-api/internal/customer/repository"
	"clientes-api/internal/customer/service"
	"clientes-api/internal/perm"
	"net/http"

	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
	"github.com/go-chi/chi/v5"
)

func Wire(lc launcher.Launcher, srv server.Server, logger logging.Logger, validator valid.Validator, db postgres.Component, authMW func(http.Handler) http.Handler, enricher authmw.IdentityEnricher, permissions security.PermissionProvider) {
	CustomRepo := repository.NewCustomer(db)
	uow := postgres.NewUnitOfWork(logger, db)
	customSRV := service.NewCustomer(logger, CustomRepo, uow)
	customerH := handler.NewCustomer(logger, validator, customSRV)

	lc.BeforeStart(func() error {
		srv.Route("/api/v1/customers", func(r chi.Router) {
			// Pública: sin AuthMiddleware, sin Enrichment, sin Authz.
			r.Post("/", customerH.Create)

			// Protegidas.
			r.Group(func(r chi.Router) {
				r.Use(authMW)                                    // 1. verifica token
				r.Use(authmw.EnrichmentMiddleware(logger, enricher)) // 2. arma identidad

				r.With(authmw.AuthzMiddleware(logger, permissions, "customers", perm.ReadCustomers)).
					Get("/", customerH.List)
				r.With(authmw.AuthzMiddleware(logger, permissions, "customers", perm.ReadCustomers)).
					Get("/{id}", customerH.Get)
				r.With(authmw.AuthzMiddleware(logger, permissions, "customers", perm.WriteCustomers)).
					Put("/{id}", customerH.Update)
				r.With(authmw.AuthzMiddleware(logger, permissions, "customers", perm.WriteCustomers)).
					Delete("/{id}", customerH.Delete)
			})
		})
		return nil
	})
}