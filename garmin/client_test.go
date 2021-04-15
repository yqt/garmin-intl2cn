package garmin

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/yqt/garmin-intl2cn/config"
	"os"
	"testing"
)

var (
	email    = config.GarminEmail
	password = config.GarminPassword
)

func TestClient_Auth(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	client := NewClient(
		Credentials(email, password),
		//SetEnv(ApiServiceHost, SsoPrefix),
		SetEnv(ApiServiceHostCn, SsoPrefixCn),
	)

	err := client.Auth(false)
	assert.Nil(t, err)
}

func TestClient_GetActivity(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	client := NewClient(
		Credentials(email, password),
		//SetEnv(ApiServiceHost, SsoPrefix),
		SetEnv(ApiServiceHostCn, SsoPrefixCn),
	)

	err := client.Auth(false)
	assert.Nil(t, err)

	//activity, err := client.GetActivity(123456)
	activity, err := client.GetActivity(123456)
	assert.Nil(t, err)
	activityBytes, err := json.Marshal(activity)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"activity": string(activityBytes),
	}).Info()
}

func TestClient_GetActivityList(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	client := NewClient(
		Credentials(email, password),
		SetEnv(ApiServiceHost, SsoPrefix),
	)

	err := client.Auth(false)
	assert.Nil(t, err)

	activityList, err := client.GetActivityList(0, 2)
	assert.Nil(t, err)

	activityListBytes, err := json.Marshal(activityList)
	assert.Nil(t, err)

	logrus.WithFields(logrus.Fields{
		"activityList": string(activityListBytes),
	}).Info()
}

func TestClient_DownloadActivity(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	client := NewClient(
		Credentials(email, password),
		SetEnv(ApiServiceHost, SsoPrefix),
	)

	err := client.Auth(false)
	assert.Nil(t, err)

	_, fileName, err := client.DownloadActivity(123456)
	assert.Nil(t, err)
	logrus.WithFields(logrus.Fields{
		"fileName": fileName,
	}).Info()
}

func TestClient_UploadActivity(t *testing.T) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	client := NewClient(
		Credentials(email, password),
		SetEnv(ApiServiceHost, SsoPrefix),
	)

	err := client.Auth(false)
	assert.Nil(t, err)

	file, fileName, err := client.DownloadActivity(123456)
	assert.Nil(t, err)

	clientCn := NewClient(
		Credentials(email, password),
		SetEnv(ApiServiceHostCn, SsoPrefixCn),
	)
	err = clientCn.Auth(false)
	assert.Nil(t, err)

	err = clientCn.UploadActivity(fileName, file)
	assert.Nil(t, err)
}
