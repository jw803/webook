package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var defaultConfiguration []byte

var env Env

type Env struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	AppRole      string `mapstructure:"APP_ROLE"`
	AppDB        string `mapstructure:"APP_DB"`
	AppRedis     string `mapstructure:"APP_REDIS"`
	APPSecretKey string `mapstructure:"APP_SECRET_KEY"`
}

func Get() *Env {
	return &env
}

func Init() {
	appEnv, ok := os.LookupEnv("APP_ENV")

	if ok && (appEnv == "staging" || appEnv == "production") {
		fmt.Printf("APP_ENV: %s\n", appEnv)
		loadFromEnvironmentVariables()
	} else {
		if appEnv == "local" {
			loadFromEnvFile(".env.json")
		} else {
			loadFromEnvFile(".env.test.json")
		}
	}
	if err := viper.Unmarshal(&env); err != nil {
		panic("environment can't be loaded")
	}
	fmt.Println("Env Loaded")
}

func loadFromEnvFile(fileName string) {
	viper.SetConfigFile(fileName)
	err := viper.ReadInConfig()
	if err != nil {
		panic("Please Download .env.json")
	}
}

func loadFromEnvironmentVariables() {
	_ = viper.BindEnv("APP_ENV")
	_ = viper.BindEnv("APP_ROLE")
	_ = viper.BindEnv("APP_DB")
	_ = viper.BindEnv("APP_REDIS`")
	_ = viper.BindEnv("APP_SECRET_KEY`")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		panic(err)
	}
}

func SetTestConfig() {
	env = Env{
		AppEnv:       "local",
		AppRole:      "webapi",
		AppDB:        "root:root@tcp(localhost:13316)/webook",
		AppRedis:     "localhost:6379",
		APPSecretKey: "",
	}
}
