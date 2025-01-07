package db

import (
	"fmt"

	_ "github.com/lib/pq"
)

type AuthRepo struct {
	Connection *Connection
}

type RandomData struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func (authRepo *AuthRepo) GetSomeRandomData() (data []RandomData, err error) {
	query, err := authRepo.Connection.DB.Query("select * from testData")
	if err != nil {
		return nil, fmt.Errorf("db.auth: query failed, %v", err)
	}
	defer func() {
		query.Close()
	}()
	var rows []RandomData
	for query.Next() {
		var row RandomData
		err = query.Scan(&row.Id, &row.Username)
		if err != nil {
			return nil, fmt.Errorf("db.auth: scanning failed, %v", err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

/**
* Test function re-creates the table and addes 2 default rows
 */
func (authRepo *AuthRepo) Seed() error {
	query, err := authRepo.Connection.DB.Query(`
	DROP TABLE IF EXISTS testData;
	
	CREATE TABLE IF NOT EXISTS testData (
		id SERIAL primary key,
		username varchar(50) unique not null
	);

	insert into testData
		(username) 
	values 
		('hello'),('world');
	`)
	if err != nil {
		return fmt.Errorf("db.auth: Query failed, %v", err)
	}
	defer func() {
		query.Close()
	}()
	return nil
}
