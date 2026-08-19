package util

import "github.com/google/uuid"


//NewID regresa un UUID v7 en string. V7 lleva  timestamp en el prefijo, asi que los IDs quedan ordenados
//por creacion - cabien para postgres
func NewId() string {
	id,err:=uuid.NewV7()

	if err!=nil{
		return uuid.NewString()
	}
	return id.String()
}