package garmin

type ActivityListItem struct {
	ActivityId     int64  `json:"activityId"`
	ActivityName   string `json:"activityName"`
	StartTimeLocal string `json:"startTimeLocal"`
	StartTimeGMT   string `json:"startTimeGMT"`
}

func (a *ActivityListItem) Equals(obj interface{}) bool {
	if a == obj {
		return true
	}

	obj1, ok := obj.(ActivityListItem)
	if !ok {
		return false
	}

	if a.ActivityId == obj1.ActivityId {
		return true
	}

	if a.StartTimeGMT == obj1.StartTimeGMT {
		return true
	}
	if a.StartTimeLocal == obj1.StartTimeLocal {
		return true
	}

	return false
}
