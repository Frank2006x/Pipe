package config

import "github.com/spf13/viper"

type Config struct {
	POSTGRES_DB string `mapstructure:"POSTGRES_DB"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	var config Config

	err = viper.Unmarshal(&config)

	return config, err

}
