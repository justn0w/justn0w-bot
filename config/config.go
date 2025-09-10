package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func LoadYmlConfig() {
	viper.SetConfigName("config.dev")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Error loading config.dev.yml file")
	}

}

func Init() {
	LoadYmlConfig()
	LoadEnv()
	//err := InitRedis()
	//if err != nil {
	//	log.Fatalf("failed to init redis index: %v", err)
	//}
	InitMySQL()
	InitMilvus()
}
