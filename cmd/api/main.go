package main

import (
	_ "github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/api"
	redisPkg "github.com/ThienKim52/golang-dev/pkg/redis"
)

// @title Health Check API
// @version 1.0
// @description This is a health check API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @schemes http
func main() {
	config, err := api.NewConfig()

	if err != nil {
		panic(err)
	}
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	app := api.NewEngine(config, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
