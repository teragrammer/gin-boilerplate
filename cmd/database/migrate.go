package main

import (
	"gin-boilerplate/database"
	"os"
)

func main() {
	args := os.Args

	database.RunMigrations()

	// run seeder
	if len(args) > 1 && args[1] == "-seed" {
		database.RunSeeder()
	}
}
