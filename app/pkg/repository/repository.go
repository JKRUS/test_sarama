package repository

import (
	"github.com/joho/godotenv"
	"github.com/kylelemons/go-gypsy/yaml"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	dataTable = "test_data"
)

type сonfig struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLMode  string
}

func newConfig() (cfg *сonfig, err error) {
	config, err := readConfigFile("configs/config.yaml")
	if err != nil {
		logrus.Fatalf("error read config %s", err.Error())
	}

	cfg, err = initDB(*config, "configs/.env")
	if err != nil {
		logrus.Fatalf("error read config key %s", err.Error())
	}
	return
}

func readConfigFile(path string) (*yaml.File, error) {
	config, err := yaml.ReadFile(path)
	if err != nil {
	}
	_, err = config.Get("kafka.ip")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("kafka.port")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("kafka.topic")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.host")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.port")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.username")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.dbname")
	if err != nil {
		return nil, err
	}
	_, err = config.Get("db.sslmode")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func initDB(c yaml.File, path string) (*сonfig, error) {
	db := &сonfig{}
	var err error
	db.Host, err = c.Get("db.host")
	if err != nil {
		return nil, err
	}
	db.Port, err = c.Get("db.port")
	if err != nil {
		return nil, err
	}
	db.UserName, err = c.Get("db.username")
	if err != nil {
		return nil, err
	}
	db.DBName, err = c.Get("db.dbname")
	if err != nil {
		return nil, err
	}
	db.SSLMode, err = c.Get("db.sslmode")
	if err != nil {
		return nil, err
	}
	if err = godotenv.Load(path); err != nil {
		return nil, err
	}
	db.Password = os.Getenv("DB_PASSWORD")
	return db, nil
}
