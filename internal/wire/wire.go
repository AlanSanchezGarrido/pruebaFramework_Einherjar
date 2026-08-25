//package wire arma la aplicacion . Es el unico lugar donde se construyen componentes y  se conectan por capas - ni main ni los
//handlers hacen new de nada

package wire

import (
	"clientes-api/internal/config"
	"clientes-api/internal/identity"
	"clientes-api/internal/repository"
	"net/http"
	"strings"

	authjwt "code.nochebuena.dev/einherjar/auth-jwt"
	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/auth/rbac"

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
	logger := logz.New(logCfg) //la configuracion que se creo del log se le pasa al constructor logz.New()

	
	db := postgres.New(logger, cfg.Postgres) 

	validator := valid.New(valid.WithMessageProvider(valid.SpanishMessages)) 

	//Evalua la API se esta ejecutando en la computadora de manera local,
	var corsMW func(http.Handler) http.Handler

	if strings.EqualFold(cfg.AppEnv, "local") {
		corsMW = mw.CORSAllowAll() //(Añade la cabecera 'Access-Control-Allow-Origin: *')
	} else {
		corsMW = mw.CORS(cfg.Server.CORSOrigins) 
	}

	customRepo := repository.NewCustomer(db)//linea encargada de hacer las consultas sql a la tabla clientes 

	enricher := identity.NewEnricher(customRepo)//aqui le estoy fabricando la heraamienta que buscara los datos de los
                                               // usuarios en la base de datos cuando se presente su token JWT
	                                           

	permissions := rbac.NewClaimsPermissionProvider("perms", authmw.GetClaims)//esta es una herramienta que se crea 
                                                                     //y nos permite verificar que permisos tiene el tocken que llego

	signer := authjwt.NewHMACSigner([]byte(cfg.JWT.Secret))//aqui le estoy diciendo todos los tockens que se creen tiene 
                                                    //tiene que llevar mi clase secreta que esta en JWT.secret de la configuracion
                                      //se utiliza una funcion constructiora que implementa la interfaz de firmado utilizando estandar HMAC
                                            
	jwtCfg := authjwt.TokenConfig{ //este es un adaptador que recibe las variables de el tiempo de vida y el emisor y las pasa 
                                 // a estre estruct del paquete autjwt para configurar un tocken y fncione 
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
		Issuer:     cfg.JWT.Issuer,
	}


	srv := server.New(logger, cfg.Server, 
		server.WithMiddleware( 
			mw.Recover(logger), 
			mw.RequestID(newRequestID),
			mw.RequestLogger(logger), 
			corsMW, //cuarto guardian: el middlewares de CORS que se configuro anterriormente
			
		
		),
	)
	//este blo que se encarga de administrarar el encendido ordenado y el apagado seguro de la base de datos y del servidor web
	//mas adelnate cuando se ejecute run prendera primero la base dedatos y despues el servidor web,
	//y cuando se apague primero apaga el service y despues la base dedatos asi evitar datos corruptos en la bd y no afecte
	lc := launcher.New(logger, cfg.Launcher) //a la funcion constructura del paquete launcher de einherjar recibe la bitacora, y configuracion
	//que define cuantos segundos esperar para que se apague
	lc.Append(db, srv) //se agregan al metodo lc que se creo, le agregamos la conexcion implementa metodos, para arrancar y cerar la bd
	//se agrega igual el srv, servidoe web (implementa métodos para escuchar peticiones y apagarse limpiamente).

	//ruta de diagnostico que vigilara si mi servidor y mi base de datos estan saludables
	//ejemplo si mi base de datos se cae mi servidor web seguira encendido recibiendo peticiones pero las peticiones mandaran un error 500 que no
	//entonces no sabremos hasta que se quejen , entonces con esto se cre una ruta para comprobar en cualquier momento si la APi y postgres esta funcionandocorrectamente
	withHealth(lc, srv, logger, cfg.Health, db)

	withCustomer(lc, srv, logger, validator, db,customRepo,permissions,signer, jwtCfg,enricher)

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
