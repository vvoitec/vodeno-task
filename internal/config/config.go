package config

import (
	"fmt"
	"os"
)

type (
	App struct {
		Env         string
		WebApiPort  string
		DatabaseURL string
		SMTPURL     string
	}
)

func MustProvide() App {
	return App{
		Env:         mustGetEnvVariable("ENV"),
		WebApiPort:  mustGetEnvVariable("WEB_API_PORT"),
		DatabaseURL: mustGetEnvVariable("DATABASE_URL"),
		SMTPURL:     mustGetEnvVariable("SMTP_URL"),
	}
}

func mustGetEnvVariable(name string) (val string) {
	val = os.Getenv(name)
	if val == "" {
		panic(fmt.Sprintf("env variable: %s is not defined", name))
	}
	return val
}
