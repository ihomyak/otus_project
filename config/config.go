package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Server struct {
		Port              string
		ReadHeaderTimeout int
		WriteTimeout      int
	}
	Database struct {
		DSN string
	}
	Redis struct {
		Addr     string
		Database int
		Password string
	}
	Logger struct {
		Level string
	}
}

func GetConfig() (*AppConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	appConfig := new(AppConfig)

	if err := viper.ReadInConfig(); err != nil {
		var _t0 viper.ConfigFileNotFoundError
		if ok := errors.Is(err, _t0); ok {
			log.Println("Warning: Config file not found. Using defaults...", err)
		} else {
			log.Println("Error during processing config", err)
			return appConfig, err
		}
	}

	appConfig.Logger.Level = viper.GetString("LOG_LEVEL")

	appConfig.Server.Port = viper.GetString("SERVER_PORT")
	appConfig.Server.ReadHeaderTimeout = viper.GetInt("SERVER_READ_HEADER_TIMEOUT")
	appConfig.Server.WriteTimeout = viper.GetInt("SERVER_WRITE_TIMEOUT")
	appConfig.Database.DSN = viper.GetString("DATABASE")
	appConfig.Redis.Addr = viper.GetString("REDIS_HOST") + ":" + viper.GetString("REDIS_PORT")
	appConfig.Redis.Database = viper.GetInt("REDIS_DATABASE")
	appConfig.Redis.Password = viper.GetString("REDIS_PASSWORD")

	return appConfig, nil
}
