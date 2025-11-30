package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	// The below package is being used indirectly
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	dbname := os.Getenv("DB_NAME")

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		return err
	}

	_, err = db.Exec("USE " + dbname)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS teachers(
			id INT AUTO_INCREMENT PRIMARY KEY, 
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			class VARCHAR(255) NOT NULL, 
			subject VARCHAR(255) NOT NULL,
			INDEX(email)
		);
	`)
	if err != nil {
		return err
	}

	return nil

}

func ConnectDB() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to mariadb!")
	return db, err
}
