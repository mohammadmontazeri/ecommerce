package helper

import (
	"fmt"
	"main/model"
)

func GetModelIfExisted(table string, id int) error {
	db := model.ConnectToDb() ;
	queryString := fmt.Sprintf("SELECT id FROM %s WHERE id=%d ",table , id)
	error := db.QueryRow(queryString).Scan(&id)

	return error
}