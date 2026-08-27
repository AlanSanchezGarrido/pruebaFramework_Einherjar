// Package handler es la capa de transporte: traduce HTTP a llamadas de service y
// de regreso. Aquí vive lo único que sabe qué es un status code o un JSON body.
//

//Es el equivalnete al @RestController de sprint: una funcion por endpoint , cero reglas de negocio dentro. Recibe dto, lo pasa model
//con los mappers de dto, llama alservice y devueve dto

package handler

import (
	"clientes-api/internal/customer/dto"
	"clientes-api/internal/customer/service"
	"context"
	"encoding/json"
	"net/http"

	"code.nochebuena.dev/einherjar/contracts/logging"
	"code.nochebuena.dev/einherjar/core/valid"
	"code.nochebuena.dev/einherjar/core/xerrors"
	"code.nochebuena.dev/einherjar/web/httputil"
	"github.com/go-chi/chi/v5"
)

//interface que contiene todos los metodos del contrato
//entonces todos los que quieran pertenecer a esta interface necesitan contener o utilizar estos metodos
type Customer interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

//struct que contiene las herramientas que necesita para atender peticiones HTTP
//entonces todos los metodos de este handler (Create, Update, Get, List, Delete) tienen acceso directo 
//a estas herramientas a travez del receptor (c *customer) osea pointer receiver
type customer struct {
	logger    logging.Logger //registro de eventos Una libreta de registros para anotar qué operaciones se hicieron (logger).
	validator valid.Validator //validador de datos //Campo que guarda el motor de validación para revisar si los datos recibidos cumplen las reglas del DTO.
	svc       service.Customer //conexcion con la logica de negocios //Campo que guarda la interfaz del Service, 
	//						permitiéndole al handler comunicarse con la lógica de negocio sin saber nada de bases de datos directas.
}

//variable descartable
var _ Customer = (*customer)(nil) 

//funcion constructora  recibe las herramientas que requieren dependencias externas para trabajar  
//es la función fábrica que recibe todas las dependencias obligatorias, las empaqueta dentro del struct 
// privado customer y devuelve la interfaz pública lista y segura para operar.
func NewCustomer(logger logging.Logger, validator valid.Validator, svc service.Customer) Customer {
	return &customer{logger: logger, validator: validator, svc: svc}
}

// Create utiliza httputil.Handle como adaptador para manejar automáticamente la petición HTTP. Primero toma el JSON que llega en el body y lo convierte
// en un dto.CreateTodo. Después valida ese DTO utilizando las validaciones definidas en sus tags. Si todo es correcto, ejecuta la función que recibe
// el DTO, llama al servicio y finalmente convierte el resultado en una respuesta JSON.

//importante
//httputil.Handle no regresa automaticamente con HHTP 200 ok cuando la funcion  termina correctamente

func (h *customer) Create(w http.ResponseWriter, r *http.Request) {

	httputil.Handle(h.validator, h.logger, func(ctx context.Context, in dto.CreateCustomer) (dto.Customer, error) {
	
		c, err := h.svc.Create(ctx, in.Tomodel(), in.Password)
		
		if err != nil {
			
			return dto.Customer{}, err
		}
		
		return dto.FromCustomer(c), nil
	})(w, r)
}
// List usa HandleNoBody: sin body de entrada, no hay nada que validar.
//
// Regresa dto.customerlist, no []dto.customer: la respuesta es un objeto con la llave
// todos adentro, nunca un arreglo pelón en la raíz. El porqué está en dto/customer.go.

func (h *customer) List(w http.ResponseWriter, r *http.Request) {
	httputil.HandleNoBody(h.logger, func(ctx context.Context) (dto.CustomerList, error) {
		cs, err := h.svc.List(ctx)
		if err != nil {
			return dto.CustomerList{}, err
		}
		return dto.FromCustomers(cs), nil
	})(w, r)
}

// Get ya no cabe en un adaptador genérico: el id viene en la ruta, no en el body.
// Ese es el escape hatch — httputil.Error y httputil.JSON dan el mismo formato de
// respuesta y el mismo logging que los adaptadores, solo que a mano.
func (h *customer) Get(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		// Nunca inventen el status: el error ya trae su código (ErrNotFound → 404)
		// y httputil.Error lo traduce y lo loguea una sola vez, al nivel correcto.
		httputil.Error(h.logger, w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FromCustomer(c))
}
// Update mezcla las dos cosas: id en la ruta y body en el request, así que hace a
// mano lo que httputil.Handle haría solo — decodificar, validar, llamar, encodear.
// Aquí sí se elige el status, porque httputil.JSON lo recibe como parámetro.
func (h *customer) Update(w http.ResponseWriter, r *http.Request) {
	
	var in dto.UpdatedCustomer
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		
		httputil.Error(h.logger, w, r, xerrors.Wrap(xerrors.ErrInvalidInput, "JSON invalid", err))
		return
	}
	if err := h.validator.Struct(in); err != nil {
		httputil.Error(h.logger, w, r, err)
		return
	}
	
	c, err := h.svc.Update(r.Context(), in.Tomodel(chi.URLParam(r, "id")))

	if err != nil {
		httputil.Error(h.logger, w, r, err)
		return
	}
	
	httputil.JSON(w,http.StatusOK,dto.FromCustomer(c))
}

//funcion Delete con PointerReceiver que pertenece al struct custome 
func (h *customer) Delete(w http.ResponseWriter, r *http.Request) {
	//en este caso como no hay un cuerpo por peticion json, si no que tenemos que sacar el id de la url de la peticion
	//entonces llamaos al service y ala funcion delete que recibe el contexto y un id ese lo sacamos de la ruta  
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		//si hubo un error se guarda en la variable err y se manda el logger osea un tipo de indicaciones por consola y termina la ejecucion
		httputil.Error(h.logger, w, r, err)
		return
	}
	//si no hubo error se muestra un status 204 que signifa que realizo un cambio en la base de datos 
	httputil.NoContent(w)
}
