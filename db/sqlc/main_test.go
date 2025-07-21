package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joekings2k/gobank/util"
	_ "github.com/lib/pq"
)



var testQueries *Queries
var testDB *sql.DB
func TestMain (m *testing.M){
	config,err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load env ")
	}
	testDB,err = sql.Open(config.DBDriver,config.DBSource)
	if err != nil{
		log.Fatal("cannot connect to db:",err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
