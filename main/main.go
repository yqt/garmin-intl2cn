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
