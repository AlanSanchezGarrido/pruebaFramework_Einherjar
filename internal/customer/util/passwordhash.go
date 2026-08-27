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
