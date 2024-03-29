package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetPath2Itself() (string, error) {
	dir, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return dir, nil
}

func BaseName(absPath string) string {
	base := filepath.Base(absPath)
	return base
}

func GetRemotePath(username, localAbsPath string) string {
	base := filepath.Base(localAbsPath)

	remoteHomePath := fmt.Sprintf("/home/%s/%s", username, base)
	return remoteHomePath
}
