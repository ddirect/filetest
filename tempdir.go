package filetest

import (
	"os"
	"testing"
)

func TempDir(t *testing.T, altName string) string {
	if os.Getenv("LOCALTEMP") == "1" {
		os.MkdirAll(altName, 0766)
		return altName
	} else {
		return t.TempDir()
	}
}
