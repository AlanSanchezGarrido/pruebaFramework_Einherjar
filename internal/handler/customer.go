//handler es la cpa de transporte: traduce HTTP allamadasde service y de regreso, aqui cice lo unico que sabe que es un
// statuscode o un JSON bosy

//Es el equivalnete al @RestController de sprint: una funcion por endpoint , cero reglas de negocio dentro. Recibe dto, lo pasa model
//con los mappers de dto, llama alservice y devueve dto

package handler

import (
	"clientes-api/internal/dto"
	"clientes-api/internal/service"
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
	//							permitiéndole al handler comunicarse con la lógica de negocio sin saber nada de bases de datos directas.
}

//variable descartable
//esta variable es descartable ya que nunca se declara por eso el (_) //nos sirvr para avisarnos si algun metodo de la interfaz nos falto por eso
//es de tipo customer que es un puntero de valor nulo

//tambien nos avisa si no esta correcto los parametros y returns de la funcion o si no coinciden con los que declaramos en la interface  

var _ Customer = (*customer)(nil) 

//funcion constructora  recibe las herramientas que requieren dependencias externas para trabajar  
//es la función fábrica que recibe todas las dependencias obligatorias, las empaqueta dentro del struct 
// privado customer y devuelve la interfaz pública lista y segura para operar.
func NewCustomer(logger logging.Logger, validator valid.Validator, svc service.Customer) Customer {
	return &customer{logger: logger, validator: validator, svc: svc}
}

//Funcion que crea un customer en este caso recibimos un cuerpo (body) que necesitamops verificar que venga bien

func (h *customer) Create(w http.ResponseWriter, r *http.Request) {
	//con el http.handle internamente verifica, valida si esta correcto y cumple con la estructura requerida para un cliente 
	//dentro de esta se crea una funcion anonima sin nombre donde se le pasa el contexto, y el objeto de transferencia y asu vez
	//returnamos el objeto o un error 
	httputil.Handle(h.validator, h.logger, func(ctx context.Context, in dto.CreateCustomer) (dto.Customer, error) {
		//aqui le pasamos a la funcion create del service lo que seria el contexto y in convertido en model.customer
		c, err := h.svc.Create(ctx, in.Tomodel())
		//si al recibir esos parametros la funcion create realiza lo que tiene que hacer y devuelve un err en la variable err se manipula 
		if err != nil {
			//si hubo error se regresa el dto.customer el struct  vacio y el error que se encontro
			return dto.Customer{}, err
		}
		//si no hubo error se regresa la funcion FromCustomer del dto, con los datos que regreso la funcion create 
		return dto.FromCustomer(c), nil
	})(w, r)//aqui hace funcionar la funcion que creo handle y le pasamos los parametros que necesita 
}

func (h *customer) List(w http.ResponseWriter, r *http.Request) {
	httputil.HandleNoBody(h.logger, func(ctx context.Context) (dto.CustomerList, error) {
		cs, err := h.svc.List(ctx)
		if err != nil {
			return dto.CustomerList{}, err
		}
		return dto.FromCustomers(cs), nil
	})(w, r)
}

func (h *customer) Get(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		httputil.Error(h.logger, w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FromCustomer(c))
}

func (h *customer) Update(w http.ResponseWriter, r *http.Request) {
	//creamos la variable donde guardara el JSON de la peticion que sea valido
	var in dto.UpdatedCustomer
		//aqui verificamos que ese json este bien con el json.NewDecoder pasandole el body(cuerpoJson) despues y se 
	// lo paso a la direccion de variable 
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		//si al momento de verificar el json esta mal cachamos el error en la variable err y despues lo mostramos 
		// con un mensaje todo esto con httputil. Error que es del framework
		httputil.Error(h.logger, w, r, xerrors.Wrap(xerrors.ErrInvalidInput, "JSON invalid", err))
		return
	}
	//si no hubo un error, validamos el json con el struct que esta en dto, si en dado caso le falta un dato que es validate:"required"
	//entnces cachamos el error dentro de err y lo manejamos mostrando el error 
	if err := h.validator.Struct(in); err != nil {
		httputil.Error(h.logger, w, r, err)
		return
	}
	//si todo sale bien y cumple con el cuerpo y la validacion, ahora sacamos el id del cliente de la ruta con chi.URLParam() y se lo pasa al metodo
	//in.Tomodel para crear el model.Customer  y despues esos datos se los pasamos al servicio en el metodo UPDATE que recibe 
	// un contexto y un model.Customer 
	c, err := h.svc.Update(r.Context(), in.Tomodel(chi.URLParam(r, "id")))
	//ahora ese metodo del service regresa un error si por alguna razon no logro realizar la logica de negocio y se majera en la variable err
	if err != nil {
		httputil.Error(h.logger, w, r, err)
		return
	}
	//si no hubi un error se regresa un status Ok y de returna un slices de ese cliente que se actualizo  ya que dto.FromCustomer()recibe un 
	//model customer
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
