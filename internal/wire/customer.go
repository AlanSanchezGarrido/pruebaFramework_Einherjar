package wire

import (
	"clientes-api/internal/handler"
	"clientes-api/internal/repository"
	"clientes-api/internal/service"
	"clientes-api/internal/perm"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"
	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
	"github.com/go-chi/chi/v5"
)

// funcion de montaje y cableado del modulo clientes
// withCustomer(...), piensa en: "Aquí se ensamblan todas las capas de Clientes y se matriculan sus endpoints justo a tiempo antes de abrir el servidor".
func withCustomer(lc launcher.Launcher, srv server.Server, logger logging.Logger, validator valid.Validator, db postgres.Component, customerRepo repository.Customer, permissions security.PermissionProvider, signer authjwt.Signer, jwtCfg authjwt.TokenConfig,enricher authmw.IdentityEnricher) {
	// BeforeStart corre después de OnInit (la base ya está abierta) y antes de
	// OnStart (el server todavía no escucha). Es el momento de registrar rutas.
	//entonces Beforstart le dice al launcher guarda esta funcion anonima y ejecuta unicamente despues de que la base de datos este conectada per
	//antes de encerder el servidor web
	lc.BeforeStart(func() error {
		                                       //Crea la capa de acceso a datos conectando directamente al compornente db
		svc := service.NewCustomer(logger, customerRepo, postgres.NewUnitOfWork(logger, db)) //crea el cerebro del negocio inyectando (logger, el repo)
		//crando al mismo tiempo la unidad de trabajo para transacciones atomaticas seguras
		h := handler.NewCustomer(logger, validator, svc) //Crea el controlador HTTP entregandole el validador de datos, el logger y el servicio recien ensamblado
        
		authSvc:=service.NewAuth(logger,customerRepo,signer,jwtCfg)
		authH:=handler.NewAuth(logger,validator,authSvc)
		srv.Route("/api/v1/auth", func(r chi.Router) {
			r.Post("/register", h.Create)
			r.Post("/login",authH.Login)
		})

		// server.Server trae chi.Router embebido, así que Route/Get/Post son los
		// de chi tal cual.
		srv.Route("/api/v1/customers", func(r chi.Router) {

			r.Use(authjwt.AuthMiddleware(logger,signer,nil),authmw.EnrichmentMiddleware(logger,enricher))
			r.With(authmw.AuthzMiddleware(logger, permissions, "customers",perm.ReadCustomers)).
				Get("/", h.List)
			r.With(authmw.AuthzMiddleware(logger, permissions, "customers", perm.ReadCustomers)).
				Get("/{id}", h.Get)
			r.With(authmw.AuthzMiddleware(logger,permissions,"customers",perm.WriteCustomers)).
			Put("/{id}", h.Update)
			r.With(authmw.AuthzMiddleware(logger,permissions,"customers",perm.WriteCustomers)).
			Delete("/{id}", h.Delete)
		})
		return nil
	})
}
