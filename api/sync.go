package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yqt/garmin-intl2cn/config"
	"github.com/yqt/garmin-intl2cn/garmin"
	"github.com/yqt/garmin-intl2cn/sync"
	"net/http"
)

func genSyncHandler(c *gin.Context) {
	userInfo := sync.UserInfo{
		Intl: garmin.UserInfo{
			Email:    config.GarminEmail,
			Password: config.GarminPassword,
		},
		Cn: garmin.UserInfo{
			Email:    config.GarminCnEmail,
			Password: config.GarminCnPassword,
		},
	}
	suc, msg, err := sync.SynchronizeLatestActivities(userInfo)

	logrus.WithFields(logrus.Fields{
		"suc": suc,
		"msg": msg,
		"err": err,
	}).Info("sync result")

	c.PureJSON(http.StatusOK, gin.H{
		"success": suc,
		"message": msg,
		"error":   fmt.Sprintf("%v", err),
	})
}
