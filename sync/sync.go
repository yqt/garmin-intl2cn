package sync

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yqt/garmin-intl2cn/garmin"
)

type UserInfo struct {
	Intl garmin.UserInfo `json:"intl"`
	Cn   garmin.UserInfo `json:"cn"`
}

func SynchronizeLatestActivities(userInfo UserInfo) (bool, string, error) {
	clientIntl := garmin.NewClient(
		garmin.Credentials(userInfo.Intl.Email, userInfo.Intl.Password),
		garmin.SetEnv(garmin.ApiServiceHost, garmin.SsoPrefix),
	)
	err := clientIntl.Auth(false)
	if err != nil {
		return false, "", err
	}

	clientCn := garmin.NewClient(
		garmin.Credentials(userInfo.Cn.Email, userInfo.Cn.Password),
		garmin.SetEnv(garmin.ApiServiceHostCn, garmin.SsoPrefixCn),
	)
	err = clientCn.Auth(false)
	if err != nil {
		return false, "", err
	}

	intlActivityList, err := clientIntl.GetActivityList(0, 3)
	if err != nil {
		return false, "", err
	}

	cnActivityList, err := clientCn.GetActivityList(0, 10)
	if err != nil {
		return false, "", err
	}

	succeedActivityIds := make([]int64, 0)
	failedActivityIds := make([]int64, 0)
	skippedActivityIds := make([]int64, 0)

	for _, intlAct := range intlActivityList {
		found := false
		for _, cnAct := range cnActivityList {
			if intlAct.Equals(cnAct) {
				found = true
				break
			}
		}
		if !found {
			file, fileName, err := clientIntl.DownloadActivity(intlAct.ActivityId)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"activityId": intlAct.ActivityId,
					"err":        err,
				}).Error("activity download failed")
				failedActivityIds = append(failedActivityIds, intlAct.ActivityId)
				continue
			}
			err = clientCn.UploadActivity(fileName, file)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"activityId": intlAct.ActivityId,
					"err":        err,
				}).Error("activity upload failed")
				failedActivityIds = append(failedActivityIds, intlAct.ActivityId)
				continue
			}
			succeedActivityIds = append(succeedActivityIds, intlAct.ActivityId)
		} else {
			skippedActivityIds = append(skippedActivityIds, intlAct.ActivityId)
		}
	}

	logrus.WithFields(logrus.Fields{
		"succeedActivityIds": succeedActivityIds,
		"failedActivityIds":  failedActivityIds,
		"skippedActivityIds": skippedActivityIds,
		"err":                err,
	}).Debug("sync detail")

	suc := true
	if len(succeedActivityIds) == 0 && len(failedActivityIds) != 0 {
		suc = false
	}
	return suc, fmt.Sprintf(
		"id[%v] succeeded. id[%v] failed. id[%v] skipped.",
		succeedActivityIds, failedActivityIds, skippedActivityIds), nil
}
