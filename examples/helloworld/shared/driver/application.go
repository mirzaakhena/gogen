package driver

import "time"

type Controller interface {
	RegisterRouter()
}

type RegistryContract interface {
	RunApplication()
}

func Run(rv RegistryContract) {
	if rv != nil {
		rv.RunApplication()
	}
}

type ApplicationData struct {
	AppName       string `json:"appName"`
	AppInstanceID string `json:"appInstanceID"`
	StartTime     string `json:"startTime"`
}

func NewApplicationData(appName, appInstanceID string) ApplicationData {
	return ApplicationData{
		AppName:       appName,
		AppInstanceID: appInstanceID,
		StartTime:     time.Now().Format("2006-01-02 15:04:05"),
	}
}
