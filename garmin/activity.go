package garmin

type Activity struct {
	ActivityId         int64        `json:"activityId"`
	ActivityName       string       `json:"activityName"`
	Description        string       `json:"description"`
	UserProfileId      int          `json:"userProfileId"`
	IsMultiSportParent bool         `json:"isMultiSportParent"`
	ActivityType       ActivityType `json:"activityTypeDTO"`
	Summary            Summary      `json:"summaryDTO"`
	MetadataDTO        MetaData     `json:"metadataDTO"`
}

type ActivityType struct {
	TypeId       int    `json:"typeId"`
	TypeKey      string `json:"typeKey"`
	ParentTypeId int    `json:"parentTypeId"`
	SortOrder    int    `json:"sortOrder"`
	IsHidden     bool   `json:"isHidden"`
}

type Summary struct {
	StartTimeLocal     string  `json:"startTimeLocal"`
	StartTimeGMT       string  `json:"startTimeGMT"`
	StartLatitude      float64 `json:"startLatitude"`
	StartLongitude     float64 `json:"startLongitude"`
	Distance           float64 `json:"distance"`
	Duration           float64 `json:"duration"`
	MovingDuration     float64 `json:"movingDuration"`
	ElapsedDuration    float64 `json:"elapsedDuration"`
	ElevationGain      float64 `json:"elevationGain"`
	ElevationLoss      float64 `json:"elevationLoss"`
	MaxElevation       float64 `json:"maxElevation"`
	MinElevation       float64 `json:"minElevation"`
	AverageSpeed       float64 `json:"averageSpeed"`
	AverageMovingSpeed float64 `json:"averageMovingSpeed"`
	MaxSpeed           float64 `json:"maxSpeed"`
	Calories           float64 `json:"calories"`
	AverageHR          float64 `json:"averageHR"`
	MaxHR              float64 `json:"maxHR"`
	AverageTemperature float64 `json:"averageTemperature"`
	MinTemperature     float64 `json:"minTemperature"`
	MaxTemperature     float64 `json:"maxTemperature"`
	TrainingEffect     float64 `json:"trainingEffect"`
	AverageBikeCadence float64 `json:"averageBikeCadence"`
	MaxBikeCadence     float64 `json:"maxBikeCadence"`
	AveragePower       float64 `json:"averagePower"`
	MaxPower           float64 `json:"maxPower"`
	AverageRunCadence  float64 `json:"averageRunCadence"`
	MaxRunCadence      float64 `json:"maxRunCadence"`
}

type MetaData struct {
	Manufacturer string   `json:"manufacturer"`
	Sensors      []Sensor `json:"sensors"`
}

type Sensor struct {
	SKU             string  `json:"sku"`
	SoftwareVersion float32 `json:"softwareVersion"`
	LocalDeviceType string  `json:"localDeviceType"`
}
