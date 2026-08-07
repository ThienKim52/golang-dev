package main

import (
	_ "github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/api"
	"github.com/ThienKim52/golang-dev/internal/app/model"
	// loggerPkg "github.com/ThienKim52/golang-dev/pkg/logger"
	redisPkg "github.com/ThienKim52/golang-dev/pkg/redis"
	"github.com/ThienKim52/golang-dev/pkg/sqldb"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/ThienKim52/golang-dev/pkg/common"
	"github.com/redis/go-redis/v9"
)

// @title Health Check API
// @version 1.1
// @description This is a health check API server.
// @Basepath /
func main() {
	// Init app config
	cfg := createAPIConfig()

	if cfg. InstanceID == "" {
	cfg.InstanceID = uuid.New().String()
	}
	
	// loggerPkg.SetLogLevel()

	// Init repos
	redisClient := createRedisClient()

	// Init DB
	db, err := sqldb.NewClient("")
	common.HandleError(err)

	err = db.AutoMigrate(&model.User{})
	common.HandleError(err)

	// Init api app
	a := createAPIApp(cfg, redisClient, db)

	// Start api app
	common.HandleError(a.Start())

}

// createRedisClient returns a redis client
func createRedisClient() *redis.Client { 
	redisClient, err := redisPkg.NewClient("")
	common.HandleError(err)

	return redisClient
}

// createAPIConfig creates API configuration based on environment variables
func createAPIConfig() *api.Config { 
	cfg, err := api.NewConfig()
	common.HandleError(err)

	return cfg
}

// createAPIApp creates API application based on API configuration
func createAPIApp(cfg *api.Config, redis *redis.Client, db *gorm.DB) api.Engine { 
	app := gin.New()
	a := api.New(app, cfg, redis, db)

	return a
}