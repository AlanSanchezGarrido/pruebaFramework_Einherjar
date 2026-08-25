package repository

import "clientes-api/internal/model"

func scanCompra(s scanner) (model.Shopping , error){
var c model.Shopping

if err :=  s.Scan(&c.ID,&c.Total,&c.Article,&c.Customer_id,&c.Create_At,&c.Update_At);err!=nil {
return model.Shopping{},nil	
}
return c,nil
}