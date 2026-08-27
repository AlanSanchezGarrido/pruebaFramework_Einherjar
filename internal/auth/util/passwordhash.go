package util

import "golang.org/x/crypto/bcrypt"

//Esta funcion CheckPassword recibe el hash que se guardo en la bd y la contraseña en texto plano, regres nil
//si coincide, error si no. No hay forma de "deshacer " el hash por eso siempre se compara nunca se desifra
func CheckPassword(hash, plain string) error  {
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(plain))
}
