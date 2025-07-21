package main

import (
	"database/sql"
	"log"

	"github.com/joekings2k/gobank/api"
	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/util"
	_ "github.com/lib/pq"
)


func main (){
	config,err :=util.LoadConfig(".")//path to config file
	if err != nil {
		log.Fatal("cannot log configs")
	}
	conn,err := sql.Open(config.DBDriver,config.DBSource)
	if err != nil{
		log.Fatal("cannot connect to db:",err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err= server.Start(config.ServerAddress)
	if err !=nil {
		log.Fatal("cannot start server:", err)
	}
}