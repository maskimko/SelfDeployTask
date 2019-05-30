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
	fmt.Println(dir)
	return dir, nil
}

func GetRemotePath(username, localAbsPath string) string {
	base := filepath.Base(localAbsPath)

	remoteHomePath := fmt.Sprintf("/home/%s/%s", username, base)
	return remoteHomePath
}
