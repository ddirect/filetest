package filetest

import "path/filepath"

type Entry struct {
	ParentPath string
	Name       string
}

func (e *Entry) Path() string {
	return filepath.Join(e.ParentPath, e.Name)
}
