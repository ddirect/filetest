package filetest

import (
	"os"
	"path/filepath"

	"github.com/ddirect/check"
)

type Dir struct {
	Entry
	Files map[string]*File
	Dirs  map[string]*Dir
}

func (d *Dir) EachDir(cb func(*Dir)) {
	for _, dir := range d.Dirs {
		cb(dir)
	}
}

func (d *Dir) EachDirRecursive(cb func(*Dir)) {
	d.EachDir(func(nd *Dir) {
		cb(nd)
		nd.EachDirRecursive(cb)
	})
}

func (d *Dir) EachFile(cb func(*File)) {
	for _, file := range d.Files {
		cb(file)
	}
}

func (d *Dir) EachFileRecursive(cb func(*File)) {
	d.EachFile(cb)
	d.EachDir(func(nd *Dir) {
		nd.EachFileRecursive(cb)
	})
}

// func (d *Dir) EachRecursive(fcb func(*File), dcb func(*Dir)) {
// 	d.EachFile(fcb)
// 	d.EachDir(func(nd *Dir) {
// 		nd.EachRecursive(fcb, dcb)
// 	})
// }

func NewDirMakerFactory(base string) func(*Dir) {
	return func(d *Dir) {
		check.E(os.Mkdir(filepath.Join(base, d.Path()), 0775))
	}
}

func (d *Dir) AllFilesMap() map[string]*File {
	files := make(map[string]*File)
	d.EachFileRecursive(func(f *File) {
		files[f.Path()] = f
	})
	return files
}

func (d *Dir) AllFilesSlice() []*File {
	var files []*File
	d.EachFileRecursive(func(f *File) {
		files = append(files, f)
	})
	return files
}

func (d *Dir) AllDirsMap() map[string]*Dir {
	dirs := make(map[string]*Dir)
	d.EachDirRecursive(func(dd *Dir) {
		dirs[dd.Path()] = dd
	})
	return dirs
}

func (d *Dir) AllDirsSlice() []*Dir {
	var dirs []*Dir
	d.EachDirRecursive(func(dd *Dir) {
		dirs = append(dirs, dd)
	})
	return dirs
}

func NewDirFromStorage(base string) *Dir {
	fileFactory := NewFileFromStorageFactory(base)
	var core func(entry Entry) *Dir
	core = func(entry Entry) *Dir {
		parentPath := entry.Path()
		entries, err := os.ReadDir(filepath.Join(base, parentPath))
		check.E(err)
		files := make(map[string]*File)
		dirs := make(map[string]*Dir)
		for _, e := range entries {
			name := e.Name()
			ne := Entry{parentPath, name}
			mode := e.Type()
			if mode.IsRegular() {
				files[name] = fileFactory(ne)
			} else if mode.IsDir() {
				dirs[name] = core(ne)
			}
		}
		return &Dir{entry, files, dirs}
	}
	return core(Entry{})
}
