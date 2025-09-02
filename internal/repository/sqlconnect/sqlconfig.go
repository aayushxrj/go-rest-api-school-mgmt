package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(dbname string) (*sql.DB, error) {

	connectionString := "root:gungun@tcp(localhost:3306)/" + dbname
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MariaDB successfully")
	return db, nil
}
