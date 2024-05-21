package main

import (
	"os"
	"path"

	"github.com/aptyInc/dbmigrate/cmd"
	"github.com/joho/godotenv"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Unable to get current working directory")
	}
	err = godotenv.Load(path.Join(cwd, ".env"))
	if err != nil {
		panic("unable to read env file")
	}
	cmd.Execute()
}
