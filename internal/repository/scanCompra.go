package repository

import "clientes-api/internal/model"

func scanCompra(s scanner) (model.Shopping , error){
var c model.Shopping

if err :=  s.Scan(&c.ID,&c.Total,&c.Article,&c.CustomerID,&c.CreatedAt,&c.UpdatedAt);err!=nil {
return model.Shopping{},nil	
}
return c,nil
}