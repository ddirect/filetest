package filetest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ddirect/check"
	"github.com/google/go-cmp/cmp"
)

type Dir struct {
	Entry
	Files []*File
	Dirs  []*Dir
}

func (d *Dir) Compare(o *Dir) bool {
	return cmp.Equal(d, o,
		cmp.Transformer("files_trans", func(files []*File) map[string]*File {
			m := make(map[string]*File, len(files))
			for _, f := range files {
				m[f.Name] = f
			}
			if len(m) != len(files) {
				panic("cmp files transformer: repeated names detected")
			}
			return m
		}),
		cmp.Transformer("dirs_trans", func(dirs []*Dir) map[string]*Dir {
			m := make(map[string]*Dir, len(dirs))
			for _, d := range dirs {
				m[d.Name] = d
			}
			if len(m) != len(dirs) {
				panic("cmp dirs transformer: repeated names detected")
			}
			return m
		}))
}

func (d *Dir) EachDir(cb func(*Dir)) {
	for _, dir := range d.Dirs {
		cb(dir)
	}
}

func (d *Dir) EachDirRecursive(cb func(*Dir)) (count int) {
	d.EachDir(func(nd *Dir) {
		cb(nd)
		count += nd.EachDirRecursive(cb) + 1
	})
	return
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

func (d *Dir) EachRecursive(fcb func(*File), dcb func(*Dir)) {
	d.EachFile(fcb)
	d.EachDir(func(nd *Dir) {
		dcb(nd)
		nd.EachRecursive(fcb, dcb)
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
	return NewDirFromStorageFiltered(base, nil)
}

func NewDirFromStorageFiltered(base string, fileFilter func(e Entry) bool) *Dir {
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
				if fileFilter == nil || fileFilter(ne) {
					dir.Files = append(dir.Files, fileFactory(ne))
				}
			} else if mode.IsDir() {
				dir.Dirs = append(dir.Dirs, core(ne))
			}
		}
		return dir
	}
	return core(Entry{})
}

// Remove files by pointer
func (d *Dir) RemoveFiles(files []*File) {
	remove := make(map[*File]bool, len(files))
	for _, file := range files {
		remove[file] = true
	}

	core := func(d *Dir) {
		files := d.Files
		w := 0
		for r, file := range files {
			if !remove[file] {
				if r != w {
					files[w] = file
				}
				w++
			}
		}
		d.Files = files[:w]
	}

	core(d)
	d.EachDirRecursive(core)
}

func (d *Dir) Dump(outFile string) {
	file, err := os.Create(outFile)
	check.E(err)
	defer check.DeferredE(file.Close)
	w := bufio.NewWriter(file)
	d.EachRecursive(
		func(f *File) {
			fmt.Fprintln(w, f.Path())
		},
		func(d *Dir) {
			fmt.Fprintf(w, "%s:\n", d.Path())
		},
	)
	w.Flush()
}
