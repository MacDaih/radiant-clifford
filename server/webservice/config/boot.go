package config

import (
	"log"
	"os"
)

type config struct {
	port string // PORT

	socket    string //SENSOR_PORT
	sensorKey string //KEY

	dbName string // DB_NAME
	dbHost string // DB_HOST
	dbPort string // DB_PORT
}

var cfg config

func set(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("failed to boot : %s is not set", key)
	}
	return value
}

func Boot() {
	cfg.port = set("PORT")

	cfg.socket = set("SENSOR_PORT")

	cfg.sensorKey = set("KEY")

	cfg.dbName = set("DB_NAME")
	cfg.dbHost = set("DB_HOST")
	cfg.dbPort = set("DB_PORT")
}

func GetPort() string {
	return cfg.port
}

func GetSocket() string {
	return cfg.socket
}

func GetSensorKey() string {
	return cfg.sensorKey
}

func GetDBEnv() (string, string, string) {
	return cfg.dbName, cfg.dbHost, cfg.dbPort
}
