package main

import (
	"ewallet-api/database"
	"ewallet-api/router"
)

const PORT = ":3000"

func main() {
	database.StartDB()

	router.StartServer().Run(PORT)
}
