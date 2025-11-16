package config

import (
	"fmt"
	"os"
)

type PostgresConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c *PostgresConf) DSN() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        c.User, c.Password, c.Host, c.Port, c.Name)
}

func LoadConfigDB() (string, error) {
    conf := &PostgresConf{
        Host:     os.Getenv("POSTGRES_HOST"),
        Port:     os.Getenv("POSTGRES_PORT"),
        User:     os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Name:     os.Getenv("POSTGRES_DB"),
    }
    if conf.Host == "" || conf.Port == "" || conf.User == "" || conf.Password == "" || conf.Name == "" {
        return "", fmt.Errorf("missing required DB env vars")
    }
    return conf.DSN(), nil
}