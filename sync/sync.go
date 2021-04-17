package sync

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yqt/garmin-intl2cn/garmin"
)

const (
	ActivityListWrapperTypeIntl = iota
	ActivityListWrapperTypeCn
)

type UserInfo struct {
	Intl garmin.UserInfo `json:"intl"`
	Cn   garmin.UserInfo `json:"cn"`
}

type ActivityListWrapper struct {
	ActivityList []garmin.ActivityListItem
	Type         int
}

func SynchronizeLatestActivities(userInfo UserInfo) (bool, string, error) {
	clientIntl := garmin.NewClient(
		garmin.Credentials(userInfo.Intl.Email, userInfo.Intl.Password),
		garmin.SetEnv(garmin.ApiServiceHost, garmin.SsoPrefix),
	)
	clientCn := garmin.NewClient(
		garmin.Credentials(userInfo.Cn.Email, userInfo.Cn.Password),
		garmin.SetEnv(garmin.ApiServiceHostCn, garmin.SsoPrefixCn),
	)

	errChan := make(chan error)
	defer close(errChan)

	actChan := make(chan ActivityListWrapper)
	defer close(actChan)

	go getActivityList(clientIntl, 0, 3, ActivityListWrapperTypeIntl, actChan, errChan)
	go getActivityList(clientCn, 0, 10, ActivityListWrapperTypeCn, actChan, errChan)

	var (
		intlActivityList []garmin.ActivityListItem
		cnActivityList []garmin.ActivityListItem
		err error
	)

	count := 0
	for count < 2 {
		select {
		case err := <-errChan:
			return false, "", err
		case actWrapper := <-actChan:
			switch actWrapper.Type {
			case ActivityListWrapperTypeIntl:
				intlActivityList = actWrapper.ActivityList
				count++
			case ActivityListWrapperTypeCn:
				cnActivityList = actWrapper.ActivityList
				count++
			}
		}

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

func getActivityList(client *garmin.Client, start int64, limit int64, actType int, resultChan chan<- ActivityListWrapper, errChan chan<- error) {
	err := client.Auth(false)
	if err != nil {
		errChan <- err
		return
	}
	activityList, err := client.GetActivityList(start, limit)
	if err != nil {
		errChan <- err
		return
	}
	activityListWrapper := ActivityListWrapper{
		ActivityList: activityList,
		Type: actType,
	}
	resultChan <- activityListWrapper
}