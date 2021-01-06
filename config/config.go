package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig hold the database connection values
type DBConfig struct {
	Host, Password, Name, User, Port string
}

// DB interface holds the sql connections open and close functionality
type DB interface {
	// connection open the connection with mysql
	open() *SQL
	// connect is the combination of three function which makes connection with db and returns
	// alive connection
	connect() *SQL
	// dcs returns the sql connection string
	dcs() string
	// loadEnv loads the env file data into running environment
	loadEnv() *DBConfig
	// setLimits sets the limited connection strings
	setLimits()
	// closed closes the open db connection
	closed()
}

// SQL holds the sql connection
type SQL struct {
	Alive *sql.DB
}

// loadEnv return the config of the database.
func loadEnv() *DBConfig {
	return &DBConfig{
		Host:     configString("DBHOST"),
		Password: configString("DBPASSWORD"),
		Name:     configString("DBNAME"),
		User:     configString("DBUSER"),
		Port:     configString("DBPORT"),
	}
}

// dcs generates the sql string
func (db *DBConfig) dcs() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", db.User, db.Password, db.Host, db.Name)
}

// connect return the sql connection pointer
func (db *DBConfig) open() (*SQL, error) {
	dbs, err := sql.Open("mysql", db.dcs())
	if err != nil {
		// log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	// setLimits(dbs)
	return &SQL{
		Alive: dbs,
	}, nil
}

// setLimits to run limited connection
func (db *SQL) setLimits() {
	db.Alive.SetMaxOpenConns(20)
	db.Alive.SetMaxIdleConns(20)
	db.Alive.SetConnMaxLifetime(time.Minute * 5)
}

// closed closes the open sql connection
func (db *SQL) closed() {
	db.Alive.Close()
}

func connect() (*SQL, error) {
	env := loadEnv()
	db, err := env.open()
	if err != nil {
		return nil, err
	}
	db.setLimits()
	return db, nil
}

// configString return the value of the env variable
func configString(name string) string {
	return os.Getenv(name)
}
