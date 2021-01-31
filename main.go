package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// << carregar as variaveis locais
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// >>

	app := App{}
	app.Initialize(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DBNAME_PROD"))
	app.Run(":" + os.Getenv("API_PORT"))
}
