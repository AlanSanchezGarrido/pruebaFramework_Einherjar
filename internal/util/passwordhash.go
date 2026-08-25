package util

import "golang.org/x/crypto/bcrypt"


//HasPassword genera el hash que se guarda en la base de datos jamas de guarda en texto plano si no se 
//que tenemos que hasearla
func HasPassword(plain string) (string, error) {
	hash,err:=bcrypt.GenerateFromPassword([]byte(plain),bcrypt.DefaultCost)

	if err!=nil {
		return "",err
	}
	return string(hash),nil
}

//Esta funcion CheckPassword recibe el hash que se guardo en la bd y la contraseña en texto plano, regres nil
//si coincide, error si no. No hay forma de "deshacer " el hash por eso siempre se compara nunca se desifra
func CheckPassword(hash, plain string) error  {
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(plain))
}
