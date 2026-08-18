package config


import (
	"github.com/caarlos0/env/v11"
	"code.nochebuena.dev/einherjar/core/launcher"
	"code.nochebuena.dev/einherjar/core/logz"
	postgres "code.nochebuena.dev/einherjar/db-postgres"
	"code.nochebuena.dev/einherjar/web/health"
	"code.nochebuena.dev/einherjar/web/server"

)

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

}

func Load() (Config,error)  {
	var cfg Config
	if err := env.Parse(&cfg);err!= nil {
		return Config{},err
	}
	return cfg,nil
}