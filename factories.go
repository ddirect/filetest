package filetest

import (
	"strings"

	"github.com/ddirect/xrand"
)

var LowerCaseChars, ValidChars string

type EntryFactory func(parentPath string) Entry
type FileFactory func(parentPath string) *File
type FilesFactory func(parentPath string) []*File
type DirFactory func(parentPath string, depth int) *Dir
type DirsFactory func(parentPath string, depth int) []*Dir

func init() {
	var b strings.Builder
	addRange := func(min, max byte) {
		for i := min; i <= max; i++ {
			b.WriteByte(i)
		}
	}
	addRange('a', 'z')
	LowerCaseChars = b.String()
	addRange('A', 'Z')
	addRange('0', '9')
	b.WriteString(" ()~#_-")
	ValidChars = b.String()
}

func NewRandomNameFactory(rnd *xrand.Xrand, sizeFactory func() int, charSet string) func() string {
	return func() string {
		var b strings.Builder
		size := sizeFactory()
		for i := 0; i < size; i++ {
			b.WriteByte(charSet[rnd.Intn(len(charSet))])
		}
		return b.String()
	}
}

func NewEntryFactory(nameFactory func() string) EntryFactory {
	return func(parentPath string) Entry {
		return Entry{parentPath, nameFactory()}
	}
}

func NullEntryFactory() EntryFactory {
	return NewEntryFactory(func() string {
		return ""
	})
}

func NewFileFactory(entryFactory EntryFactory) FileFactory {
	return func(parentPath string) *File {
		return &File{entryFactory(parentPath), nil}
	}
}

func NewDirFactory(entryFactory EntryFactory, filesFactory FilesFactory, dirsFactory DirsFactory) DirFactory {
	return func(parentPath string, depth int) *Dir {
		entry := entryFactory(parentPath)
		entryPath := entry.Path()
		files := make(map[string]*File)
		dirs := make(map[string]*Dir)
		for _, file := range filesFactory(entryPath) {
			files[file.Name] = file
		}
		for _, dir := range dirsFactory(entryPath, depth+1) {
			dirs[dir.Name] = dir
		}
		return &Dir{entry, files, dirs}
	}
}

func NewFilesFactory(fileFactory FileFactory, countFactory func() int) FilesFactory {
	return func(parentPath string) []*File {
		count := countFactory()
		files := make([]*File, count)
		for i := range files {
			files[i] = fileFactory(parentPath)
		}
		return files
	}
}

func NewDirsFactory(dirFactory DirFactory, countFactory func() int) DirsFactory {
	return func(parentPath string, depth int) []*Dir {
		count := countFactory()
		dirs := make([]*Dir, count)
		for i := range dirs {
			dirs[i] = dirFactory(parentPath, depth)
		}
		return dirs
	}
}
