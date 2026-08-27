package identity

import (
	"clientes-api/internal/customer/repository"
	"context"

	"code.nochebuena.dev/einherjar/auth/authmw"
	"code.nochebuena.dev/einherjar/contracts/security"
)

//estructura guarda una conexcion a mi repository  de clientes  para poder ir a buscar los datos completos del cliente
//a PostgresAQL en cuanto se valida el JWT, entonces tiene las herramientas necesarias para buscar los datos del usuairo
// autenticarlo y cargarlos en la memoria de la peticion
type Enricher struct {
	customers repository.Customer
}
//creamos una variable descartable que nos ayuda a establecer correctamente el Enricher sin que nos falte nada 
var _ authmw.IdentityEnricher = (*Enricher)(nil)

//Funcion Constructora crea las herramientas que ocuparemos 
//Función constructora que recibe el repositorio de clientes ya conectado y devuelve una 
// instancia de Enricher lista para inyectarse al middleware de autenticación.
func NewEnricher(customers repository.Customer) *Enricher  {
	return &Enricher{customers: customers}
}
//funcion Enricher  para recibir los datos id , saca el id que se empaqueto en el JWT y lo va a busacr si existe
//returna  un tipo struct que actua como gafete con el ID,NAME,EMAIL  

func (e *Enricher)Enrich(ctx context.Context, uid string, claims map[string]any ) (security.Identity, error) {

	c,err:=e.customers.Get(ctx,uid)
	
	if err!= nil {
		return security.Identity{},err
	}

	return security.Identity{UID: c.ID,DisplayName: c.Name,Email: c.Email},nil
}

