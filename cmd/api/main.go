package main

import (
	_ "github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/api"
	redisPkg "github.com/ThienKim52/golang-dev/pkg/redis"
	loggerPkg "github.com/ThienKim52/golang-dev/pkg/logger"
)

// @title Health Check API
// @version 1.1
// @description This is a health check API server.
// @Basepath /
func main() {
	config, err := api.NewConfig()
	
	//set log level
	loggerPkg.SetLogLevel(config.LogLevel)

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
