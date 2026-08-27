package wire

import (
	"clientes-api/internal/Compras/handler"
	"clientes-api/internal/Compras/repository"
	"clientes-api/internal/Compras/service"
	"net/http"

	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
	"github.com/go-chi/chi/v5"
)

func Wire(lc launcher.Launcher, srv server.Server, logger logging.Logger, validator valid.Validator, db postgres.Component, authMW func(http.Handler) http.Handler, enricher authmw.IdentityEnricher) {
	shoppingRepo := repository.NewShopping(db)
	uow := postgres.NewUnitOfWork(logger, db)
	shoppingSrv := service.NewShopping(logger, shoppingRepo, uow)
	shoppingH := handler.NewShopping(logger, validator, shoppingSrv)

	lc.BeforeStart(func() error {
		srv.Route("/api/v1/shopping", func(r chi.Router) {
			r.Use(authMW)
			r.Use(authmw.EnrichmentMiddleware(logger, enricher))

			r.Get("/", shoppingH.List)
			r.Get("/{id}", shoppingH.Get)
			r.Post("/", shoppingH.Create)
			r.Put("/{id}", shoppingH.Update)
			r.Delete("/{id}", shoppingH.Delete)
		})
		return nil
	})
}