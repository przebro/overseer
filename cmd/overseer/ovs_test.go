package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestInitialize(t *testing.T) {

	tempdir := os.TempDir()
	execPath := filepath.Join(tempdir, "bin", "overseer.exe")
	args := []string{execPath}
	fmt.Println(args)

}

func TestPushTaskToActive(t *testing.T) {

}
func AddResource(t *testing.T) {

}
func RemoveResource(t *testing.T) {

}
