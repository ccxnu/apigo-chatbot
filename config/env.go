package config

import (
	"log"

	"github.com/spf13/viper"
)

// Env now only contains Database configuration from config.json
// All other configuration is accessed directly from ParameterCache
type Env struct {
	Database DatabaseConfig `mapstructure:"Database"`
}

type DatabaseConfig struct {
	Host          string `mapstructure:"HOST"`
	Port          int    `mapstructure:"PORT"`
	User          string `mapstructure:"USER"`
	Password      string `mapstructure:"PASSWORD"`
	Name          string `mapstructure:"NAME"`
	MaxConnection int    `mapstructure:"MAX_CONNECTION"`
}

// NewEnv loads only Database configuration from config.json
func NewEnv() *Env {
	env := Env{}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error: No se puede encontrar el archivo de configuración: ", err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatal("Error: La configuración no se puede cargar en la estructura: ", err)
	}

	return &env
}
