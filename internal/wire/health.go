package wire

import (
	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/contracts/observability"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/web/health"
	"code.nochebuena.dev/einherjar/web/server"
)

// withHealth arma el endpoint de salud a partir de los componentes Checkable de
// la app (aquí nada más la base). Los checks corren en paralelo con el timeout de
// EINHERJAR_HEALTH_CHECK_TIMEOUT.
func withHealth(lc launcher.Launcher, srv server.Server, logger logging.Logger, cfg health.Config, checks ...observability.Checkable) {
	lc.BeforeStart(func() error {
		srv.Get("/health", health.NewHandlerWithConfig(logger, cfg, checks...).ServeHTTP)
		return nil
	})
}
