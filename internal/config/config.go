package config

import (
	"github.com/joho/godotenv"
	"os"
)

var Config Configuration

// Configuration settings of the service.
type Configuration struct {
	// Host --> Default 0.0.0.0
	Host string `mapstructure:"HOST" json:"host" default:"127.0.0.1"`
	// Port --> Default 49100
	Port string `mapstructure:"PORT" json:"port" default:"3100"`
	// DBName --> Name of the database. Default "amc.db"
	DBName string `mapstructure:"DB_NAME" json:"DBName" default:"amc.db"`
}

func LoadConfiguration() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = godotenv.Load(pwd + "/internal/config/.env")
	if err != nil {
		return err
	}
	Config.Host = os.Getenv("HOST")
	Config.Port = os.Getenv("PORT")
	Config.DBName = os.Getenv("DB_NAME")
	return nil
}
