package util

import "encoding/json"

func MustJSON(obj any) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
