package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig() viper.Viper {

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("database", "data.db")
	viper.SetDefault("ssh.port", "2222")
	viper.SetDefault("ssh.keypath", "ssh_private_key.pem")
	viper.SetDefault("ssh.newuser", "new")
	viper.SetDefault("admingroup", "admin")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		// create default config file
		viper.SafeWriteConfig()
	}

	return *viper.GetViper()
}
