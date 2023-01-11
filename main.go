package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"

	"gateway-api/domain"
	httpHandler "gateway-api/gateway/delivery/http"
	"gateway-api/gateway/usecase"
	"gateway-api/helpers"
	"gateway-api/helpers/response"

	log "github.com/sirupsen/logrus"
)

var goenv = ""

func init() {
	_ = godotenv.Load()

	goenv = helpers.GoEnv()
}

func main() {
	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "5"
	}

	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	timeoutContext := time.Duration(timeout) * time.Second
	logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
	if logMaxSize == 0 {
		logMaxSize = 5 //default 50 megabytes
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "server.log",
		MaxSize:    logMaxSize,
		MaxBackups: 1,
		LocalTime:  true,
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(io.MultiWriter(lumberjackLogger, os.Stdout))

	serviceConfig := readConfig()

	if goenv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = log.StandardLogger().Writer()
	router := gin.New()

	// cors
	router.Use(cors.New(serviceConfig.Cors))

	router.Use(customRequestLogger())
	router.Use(customRecovery(goenv))

	cu := usecase.NewUsecase(timeoutContext, serviceConfig)
	httpHandler.NewHandler(router, cu)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	router.Run(":" + port)
}

func customRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		log.WithFields(log.Fields{
			"ip":      c.ClientIP(),
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"proto":   c.Request.Proto,
			"status":  c.Writer.Status(),
			"latency": time.Since(startTime),
			"ua":      c.Request.UserAgent(),
		}).Info()
	}
}

func customRecovery(goenv string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Catch Error:", err)
				if goenv == "development" || goenv == "local" {
					fmt.Println(string(debug.Stack()))
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(500, "Something went wrong"))
			}
		}()
		c.Next()
	}
}

func readConfig() domain.Config {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	config := domain.Config{}
	json.Unmarshal(b, &config)

	return config
}
