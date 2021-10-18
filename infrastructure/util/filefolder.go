package util

import "os"

func IsFileExist(filepath string) bool {
  if _, err := os.Stat(filepath); os.IsNotExist(err) {
    return false
  }
  return true
}
