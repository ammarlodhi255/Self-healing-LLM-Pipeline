package utils

import "os"

// SetupTempFolders creates the temp output directory for compiled files, panics if it fails
func SetupTempFolders(tempOutputDir string) {
	// 0777 are the permissions for the directory, everyone can read, write and execute
	err := os.MkdirAll(tempOutputDir, os.ModePerm)
	if err != nil {
		panic("Error creating temp output directory:\n\n" + err.Error())
	}
}

// RemoveTempFolders removes the temp output directory for compiled files, panics if it fails
func RemoveTempFolders(tempOutputDir string) {
	err := os.RemoveAll(tempOutputDir)
	if err != nil {
		panic("Error removing temp output directory:\n\n" + err.Error())
	}
}
