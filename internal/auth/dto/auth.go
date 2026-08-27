package dto



//molde que recibe, decodifica y valida que los datos que envio un usuario sean validos 
//este contenedor lo ocupara handler al momento de verificar que el usuario mando su correo y su clave antes de procesar
//la autenticacion
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}


//Estructura que se devuelve al cliente  
// contiene la informacion del cliente y sus tokens de autenticación 
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

//Funcion encargada de recibir el modelo de datos, los tokens generados y el tiempo de expiracion y emsamblarlo todo
//en un struct llamaod AuthResponse al momento de logearse un usuario
func FromAuth( accesToken, refreshToken string, expiresIn int64) AuthResponse  {
	return AuthResponse{	
			AccessToken: accesToken,
			RefreshToken: refreshToken,
			ExpiresIn: expiresIn,
	}
}