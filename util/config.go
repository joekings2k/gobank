package util

import "github.com/spf13/viper"

// all
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}


func LoadConfig(path string) (config Config,err error){
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")//could be json or yml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil{
		return
	}
	err= viper.Unmarshal(&config)
	return
}