package config

import "github.com/spf13/viper"

type Config struct {
	POSTGRES_DB string `mapstructure:"POSTGRES_DB"`
	GITHUB_CLIENT_ID string `mapstructure:"GITHUB_CLIENT_ID"`
	GITHUB_CLIENT_SECRET string `mapstructure:"GITHUB_CLIENT_SECRET"`
	GITHUB_CALLBACK_URL string `mapstructure:"GITHUB_CALLBACK_URL"`
	JWT_SECRET string `mapstructure:"JWT_SECRET"`
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
