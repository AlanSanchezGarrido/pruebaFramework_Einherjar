package wire

import (
	shoppingwire "clientes-api/internal/Compras/wire"
	"net/http"

	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
)

// wire_shopping.go
func whithShopping(lc launcher.Launcher, srv server.Server, logger logging.Logger, validator valid.Validator, db postgres.Component, authMW func(http.Handler) http.Handler, enricher authmw.IdentityEnricher) {
	shoppingwire.Wire(lc, srv, logger, validator, db, authMW, enricher)
}