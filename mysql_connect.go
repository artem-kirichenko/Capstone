package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DB DBConfig `json:"mysql"`
}
type DBConfig struct {
	Host     string `json:host`
	Port     string `json:port`
	Username string `json:username`
	Password string `json:password`
	Database string `json:database`
}

func mysql_connect() (*sql.DB, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(filepath.Join(pwd, "env", "env.json"))
	if err != nil {
		return nil, err
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	//fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
	//	config.DB.Username,
	//	config.DB.Password,
	//	config.DB.Host,
	//	config.DB.Port,
	//	config.DB.Database))

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.DB.Username,
			config.DB.Password,
			config.DB.Host,
			config.DB.Port,
			config.DB.Database))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil

}
