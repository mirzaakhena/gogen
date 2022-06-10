package payload

import "helloworld/shared/driver"

type Args struct {
	Type      string                 `json:"type"`
	Data      any                    `json:"data"`
	Publisher driver.ApplicationData `json:"publisher"`
	TraceID   string                 `json:"traceId"`
}

type Reply struct {
	Success      bool
	ErrorMessage string
	Data         any
}
