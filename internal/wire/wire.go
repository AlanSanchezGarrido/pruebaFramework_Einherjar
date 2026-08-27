//package wire arma la aplicacion . Es el unico lugar donde se construyen componentes y  se conectan por capas - ni main ni los
//handlers hacen new de nada

package wire

import (
	"clientes-api/internal/config"

	"net/http"
	"strings"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"

	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/logz"
	"code.nochebuena.dev/einherjar/core/valid"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/mw"
	"code.nochebuena.dev/einherjar/web/server"
	"github.com/google/uuid"
)

//Run carga la configuracion, construye infraestructura, registra los hoks de cada feature y se bloquea hasta el shutdown, el orfen importa :
// config,logger, infra,launcher con todos los componentes, luego los hoks y al final lc.Run()

func Run() error {

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logCfg := cfg.Log
	logCfg.StaticArgs = []any{"service", "clientes-api", cfg.AppEnv}
	logger := logz.New(logCfg)

	//Instancia de la base de datos
	db := postgres.New(logger, cfg.Postgres)

	//Instancia del validador
	validator := valid.New(valid.WithMessageProvider(valid.SpanishMessages))



	//Evalua la API se esta ejecutando en la computadora de manera local,
	var corsMW func(http.Handler) http.Handler

	if strings.EqualFold(cfg.AppEnv, "local") {
		corsMW = mw.CORSAllowAll() //(Añade la cabecera 'Access-Control-Allow-Origin: *')
	} else {
		corsMW = mw.CORS(cfg.Server.CORSOrigins)
	}

	srv := server.New(logger, cfg.Server,
		server.WithMiddleware(
			mw.Recover(logger),
			mw.RequestID(newRequestID),
			mw.RequestLogger(logger),
			corsMW, //cuarto guardian: el middlewares de CORS que se configuro anterriormente

		),
	)

	//a la funcion constructura del paquete launcher de einherjar recibe la bitacora, y configuracion
	//que define cuantos segundos esperar para que se apague
	lc := launcher.New(logger, cfg.Launcher)
	lc.Append(db, srv)
	jwtCfg := authjwt.TokenConfig{

		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
		Issuer:     cfg.JWT.Issuer,
	}
	authMW,enricher,permissions := withAuth(lc,srv,logger,validator,db,jwtCfg,cfg.JWT.Secret)

	withHealth(lc, srv, logger, cfg.Health, db)

	withCustomer(lc, srv, logger, validator, db, authMW,enricher,permissions)

	whithShopping(lc, srv, logger, validator, db, authMW, enricher)

	return lc.Run()

}

// newRequestID regresa un UUID v7 (ordenado por tiempo), con fallback a v4.
func newRequestID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}
