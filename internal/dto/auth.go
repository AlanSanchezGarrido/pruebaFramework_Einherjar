package dto

import "clientes-api/internal/model"


//molde que recibe, decodifica y valida que los datos que envio un usuario sean validos 
//este contenedor lo ocupara handler al momento de verificar que el usuario mando su correo y su clave antes de procesar
//la autenticacion
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

//structura que responde al cliente con la credencial de acceso al sisteam, con su tiempo de renovacion, y cuando se expira en
//segundos 
type AuthTokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	ExpiresIn   int64 `json:"expires_in"`
}


//Estructura que se devuelve al cliente con el cliente 
// contiene la informacion del cliente y sus tokens de autenticación 
type AuthResponse struct {
	Customer Customer   `json:"customer"`
	Tokens   AuthTokens `json:"tokens"`
}

//Funcion encargada de recibir el modelo de datos, los tokens generados y el tiempo de expiracion y emsamblarlo todo
//en un struct llamaod AuthResponse al momento de logearse un usuario
func FromAuth(c model.Customer, accesToken, refreshToken string, expiresIn int64) AuthResponse  {
	return AuthResponse{
		Customer: FromCustomer(c),
		Tokens: AuthTokens{
			AccessToken: accesToken,
			RefreshToken: refreshToken,
			ExpiresIn: expiresIn,
		},
	}
}
