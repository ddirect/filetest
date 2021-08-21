package filetest

import "path/filepath"

type Entry struct {
	Parent *Dir
	Name   string
}

func (e *Entry) Path() string {
	if e.Parent == nil {
		return e.Name
	}
	return filepath.Join(e.Parent.Path(), e.Name)
}

func (e *Entry) PathFrom(base string) string {
	return filepath.Join(base, e.Path())
}
