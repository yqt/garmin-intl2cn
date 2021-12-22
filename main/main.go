package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yqt/garmin-intl2cn/api"
	"net/http"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	logLvlStr := os.Getenv("LOG_LEVEL")
	var logLvl logrus.Level
	switch logLvlStr {
	case "debug":
		logLvl = logrus.DebugLevel
	case "info":
		logLvl = logrus.InfoLevel
	default:
		logLvl = logrus.InfoLevel
	}
	logrus.SetLevel(logLvl)

	port := os.Getenv("PORT")

	if port == "" {
		port = "38080"
	}

	err := api.InitRoute(r)
	if err != nil {
		logrus.Fatal(err)
	}

	if err = http.ListenAndServe("localhost:"+port, r); err != nil {
		logrus.Fatal(err)
	}
}
