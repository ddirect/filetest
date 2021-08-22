package filetest

import (
	"os"
	"path/filepath"
	"testing"
)

func TempDir(t *testing.T, altName string) string {
	if temp, ok := os.LookupEnv("TESTDIR"); ok {
		res := filepath.Join(temp, altName)
		os.MkdirAll(res, 0766)
		return res
	} else {
		return t.TempDir()
	}
}
