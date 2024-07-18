package main

import (
	"fmt"
	"os"
)

func main() {
	envVars := []struct {
		name         string
		description  string
		defaultValue string
	}{
		{
			name:         "CONFIG_PATH",
			description:  "path to config file",
			defaultValue: "./config/local.yml",
		},
		{
			name:         "POSTGRES_HOST",
			description:  "Postgres database host (local: localhost, docker: db)",
			defaultValue: "db",
		},
		{
			name:         "POSTGRES_PORT",
			description:  "Postgres database port",
			defaultValue: "5432",
		},
		{
			name:         "POSTGRES_USER",
			description:  "Postgres database user",
			defaultValue: "username",
		},
		{
			name:         "POSTGRES_PASSWORD",
			description:  "Postgres database password",
			defaultValue: "password",
		},
		{
			name:         "POSTGRES_DB",
			description:  "Postgres database name",
			defaultValue: "database",
		},
	}
	file, err := os.OpenFile(".env", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = file.Close()
	}()

	for _, envVar := range envVars {
		var value string

		fmt.Printf("Insert %s (default: %s): ", envVar.description, envVar.defaultValue)
		_, err := fmt.Scanln(&value)
		if err != nil {
			if err.Error() != "unexpected newline" {
				panic(err)
			} else {
				value = ""
			}
		}

		if value == "" {
			value = envVar.defaultValue
		}

		_, err = fmt.Fprintf(file, "%s=%s\n", envVar.name, value)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Environment variables set")
}
