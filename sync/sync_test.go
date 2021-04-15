package sync

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yqt/garmin-intl2cn/config"
	"github.com/yqt/garmin-intl2cn/garmin"
	"os"
	"testing"
)

var (
	userInfo = UserInfo{
		Intl: garmin.UserInfo{
			Email:    config.GarminEmail,
			Password: config.GarminPassword,
		},
		Cn: garmin.UserInfo{
			Email:    config.GarminCnEmail,
			Password: config.GarminCnPassword,
		},
	}
)

func TestSynchronizeLatestActivities(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	suc, msg, err := SynchronizeLatestActivities(userInfo)
	assert.True(t, suc)
	logrus.WithFields(logrus.Fields{
		"msg": msg,
		"err": err,
	}).Info()
}
