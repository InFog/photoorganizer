package utils

import (
	"os"
)

func IsFile(fileName string) bool {
	isfile := true

	if fileInfo, err := os.Stat(fileName); err != nil {
		isfile = false
	} else {
		if fileInfo.IsDir() {
			isfile = false
		}
	}

	return isfile
}
