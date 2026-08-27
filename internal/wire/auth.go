package wire

import (
	auth "clientes-api/internal/auth/wire"
	"net/http"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"
	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/security"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/server"
)

func withAuth(lc launcher.Launcher,srv server.Server,logger logging.Logger,validator valid.Validator,db postgres.Component,jwtCfg authjwt.TokenConfig,jwtSecret string) (func(http.Handler) http.Handler, authmw.IdentityEnricher, security.PermissionProvider) {
return auth.Wire(lc,srv,logger,validator,db,jwtCfg,jwtSecret)
}