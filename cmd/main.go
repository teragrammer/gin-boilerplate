package main

import (
	"gin-boilerplate/pkg"
	middlewares2 "gin-boilerplate/pkg/middlewares"
	"gin-boilerplate/pkg/routes"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// generate a new log file every time the app restarts
	// log to a file
	// Get the current time
	now := time.Now()
	// Format the date and time
	formattedTime := now.Format("2006-01-02 15:04:05")
	file, err := os.OpenFile(strings.ReplaceAll(".logs/"+formattedTime, " ", "-")+".main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// load boot configuration file
	var boot = pkg.InitBoot("env.json", nil)

	// rate limit
	boot.Engine.Use(middlewares2.RateLimitMiddleware(boot))

	// CORS
	boot.Engine.Use(middlewares2.CORSMiddleware())

	// routes
	routes.V1Routes(boot)

	err = boot.Engine.Run(":" + strconv.Itoa(boot.Env.App.Port))
	if err != nil {
		return
	} // listen and serve
}
