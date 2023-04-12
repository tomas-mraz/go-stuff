package main

import (
	//	"fmt"
	"os"
	"path/filepath"
)

var dataDir string

func GetAppDir() string {
	appFile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(appFile)
}

func GetOsAppDataDir() string {
	osAppDataDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.FromSlash(osAppDataDir + "/orionx")
}

func ExistDir(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func SetDataDir(path string) {
	dataDir = path
	os.MkdirAll(dataDir, 0700)
}

func GetDataDir() string {
	return dataDir
}

func DataFile(filename string) string {
	return filepath.FromSlash(filepath.Join(dataDir, filename))
}
