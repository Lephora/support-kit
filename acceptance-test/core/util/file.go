package util

import (
	"fmt"
	"os"
	"strings"
)

func FileOrDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true

		}
		return false
	}
	return true
}

func FileExistWithExtensionName(path string, extensionName ...string) (string, bool) {
	for _, extension := range extensionName {
		fp := fmt.Sprintf("%s.%s", path, extension)
		if FileOrDirExist(fp) {
			return fp, true
		}
	}
	return "", false
}

func GetFileNameWithoutExtension(fileName string) string {
	index := strings.LastIndex(fileName, ".")
	caseName := fileName[:index]
	return caseName
}
