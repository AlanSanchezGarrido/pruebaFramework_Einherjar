package wire

import (
	customerwire "clientes-api/internal/customer/wire" // <- el paquete wire DEL MÓDULO
	"net/http"

	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
)

func withCustomer(lc launcher.Launcher,srv server.Server,logger logging.Logger, validator valid.Validator, db postgres.Component,authMW func(http.Handler) http.Handler,enricher authmw.IdentityEnricher,permissions security.PermissionProvider)  {
		customerwire.Wire(lc,srv,logger,validator,db,authMW,enricher,permissions)
}
