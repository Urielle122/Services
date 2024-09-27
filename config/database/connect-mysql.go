package database

import (
	"database/sql"
	"services/logs"

	_ "github.com/go-sql-driver/mysql"
)


func ConnectMysql()(*sql.DB,error){

	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3307)/admin" )
	if err != nil {
		logs.Errorf("Erreur lors de l'ouverture de la connexion :", err)
	}
	logs.Info("Connected to mysql")
	return db, err
}