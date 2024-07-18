package env

import "github.com/joho/godotenv"

func MustLoad() {
	if err := godotenv.Load(); err != nil {
		panic("cannot load env file: " + err.Error())
	}
}
