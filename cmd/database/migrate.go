package main

import "gin-boilerplate/database"

func main() {
	database.RunMigrations()
}
