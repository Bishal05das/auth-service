package database

import (
	"database/sql"
	_ "github.com/lib/pq"

)


type Database struct {
	DB *sql.DB
}


func NewDatabase(connectionString string) (*Database, error) {
	//open database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	//verify connection is working
	if err := db.Ping(); err != nil {
		return nil,err
		
	}
	return &Database{DB: db},nil 


}