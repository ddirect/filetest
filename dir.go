package filetest

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/ddirect/check"
)

type Dir struct {
	Entry
	Files []*File
	Dirs  []*Dir
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

func (d *Dir) Sort() {
	sort.Slice(d.Files, func(i, j int) bool {
		return d.Files[i].Name < d.Files[j].Name
	})
	sort.Slice(d.Dirs, func(i, j int) bool {
		return d.Dirs[i].Name < d.Dirs[j].Name
	})
	d.EachDirRecursive(func(nd *Dir) {
		nd.Sort()
	})
}

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
		entries, err := os.ReadDir(filepath.Join(base, entry.Path()))
		check.E(err)
		dir := &Dir{Entry: entry}
		for _, e := range entries {
			name := e.Name()
			ne := Entry{dir, name}
			mode := e.Type()
			if mode.IsRegular() {
				dir.Files = append(dir.Files, fileFactory(ne))
			} else if mode.IsDir() {
				dir.Dirs = append(dir.Dirs, core(ne))
			}
		}
		return dir
	}
	return core(Entry{})
}
