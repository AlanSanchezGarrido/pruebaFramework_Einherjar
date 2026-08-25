package main

import (
	"clientes-api/internal/wire"
	_ "github.com/joho/godotenv/autoload"
	"fmt"
	"os"
)

func main() {

	//este enciende toda la aplicacion solo eso con el metodo Run del modulo wire, entonces lo que dice llama a la funcion Run de wire 
	//enciende el servidor y si todo funciona mantente prendido pero si no apagate 

	//en pocas palabras  enciende el motor de la API y garantiza una salida controlada en caso de fallar 
	if err:=wire.Run();err!=nil {
		fmt.Fprintln(os.Stderr,"fatal:",err)
		//i ocurre un problema al arrancar, el programa no se queda trabado: reporta el error exacto en el canal adecuado (os.Stderr) y le avisa al sistema 
		// operativo con os.Exit(1) para que herramientas como Docker o Kubernetes se enteren de inmediato y reinicien el servicio.
		os.Exit(1)
	}
}