package entity

import (
	"fmt"
)

// GetLogRootFolderName ...
func GetLogRootFolderName() string {
	return fmt.Sprintf("infrastructure/log")
}

// GetLogInterfaceFileName ...
func GetLogInterfaceFileName() string {
	return fmt.Sprintf("%s/log._go", GetLogRootFolderName())
}

// GetLogImplementationFileName ...
func GetLogImplementationFileName() string {
	return fmt.Sprintf("%s/log_default._go", GetLogRootFolderName())
}
