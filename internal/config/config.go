package config

import (
	"time"

	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/logz"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/health"
	"code.nochebuena.dev/einherjar/web/server"
	"github.com/caarlos0/env/v11"
)

//JWT este struct agrupa lo que auth-JWT necesita para firmar y verificar tokens
//Ninguno de estos valores los leee el modulo por su cuento por eso estan y viven aqui

type JWT struct{
	Secret     string        `env:"JWT_SECRET,required"`//busca en el archivo .env la varibale JWT_SECRET y le asigna
														//el valor a secret, y es importante tenerla
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`//asigna un tiempo estimado de validacion de 15 minutos
																	//si en ese laop no se ocupa o se ocupo el tocken desaparese
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"168h"`//Tiempo de vida de actualizacion del tocken si no se define
																	//toma por defecto 168 horas 
	Issuer     string        `env:"JWT_ISSUER" envDefault:"clientes-api"` //Nombre asignado por defecto para firmar 
																			//el campo iss dentro de los claims
}

//config es la configuracion resuelta. caarlos0/env baja recursivamente a los config del framework 
//y rellena los tags
type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"local"`

	//config de componentes del framework -compuestos tal cual; sus tags

	Launcher launcher.Config //Einherjar_component_stop_timeout
	Log      logz.Config//EINHERJAR_LOG_*
	Server server.Config
	Health health.Config
	Postgres postgres.Config
	JWT JWT

}

func Load() (Config,error)  {
	var cfg Config
	if err := env.Parse(&cfg);err!= nil {
		return Config{},err
	}
	return cfg,nil
}