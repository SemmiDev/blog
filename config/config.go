package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var Env Configuration

type Configuration struct {
	DBDriver               string        `mapstructure:"DB_DRIVER"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	ServerAddress          string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	FirebaseCredentialJSON string        `mapstructure:"FIREBASE_CREDENTIAL_JSON"`
	FirebaseBucketName     string        `mapstructure:"FIREBASE_BUCKET_NAME"`
}

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&Env)
}
