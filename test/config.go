package main

import "github.com/jmoiron/sqlx"

type Config struct {
	Host          string
	Port          string
	User          string
	Password      string
	Db_name       string
	Max_open_cons int
	Max_idle_cons int
}

var db *sqlx.DB
var config *Config

func initConfig() {
	config.Host = "127.0.0.1"
	config.Port = "3306"
	config.User = "root"
	config.Password = "root"
	config.Db_name = "sql_test"
	config.Max_open_cons = 200
	config.Max_idle_cons = 50
}
