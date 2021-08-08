package filetest

import (
	"strings"

	"github.com/ddirect/xrand"
)

var LowerCaseChars, ValidChars string

type EntryFactory func(parent *Dir) Entry
type FileFactory func(parent *Dir) *File
type FilesFactory func(parent *Dir) []*File
type DirFactory func(parent *Dir, depth int) *Dir
type DirsFactory func(parent *Dir, depth int) []*Dir

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
	return func(parent *Dir) Entry {
		return Entry{parent, nameFactory()}
	}
}

func NullEntryFactory() EntryFactory {
	return NewEntryFactory(func() string {
		return ""
	})
}

func NewFileFactory(entryFactory EntryFactory) FileFactory {
	return func(parent *Dir) *File {
		return &File{entryFactory(parent), nil}
	}
}

func NewDirFactory(entryFactory EntryFactory, filesFactory FilesFactory, dirsFactory DirsFactory) DirFactory {
	return func(parent *Dir, depth int) *Dir {
		dir := &Dir{Entry: entryFactory(parent)}
		dir.Files = filesFactory(dir)
		dir.Dirs = dirsFactory(dir, depth+1)
		return dir
	}
}

func NewFilesFactory(fileFactory FileFactory, countFactory func() int) FilesFactory {
	return func(parent *Dir) []*File {
		count := countFactory()
		files := make([]*File, count)
		for i := range files {
			files[i] = fileFactory(parent)
		}
		return files
	}
}

func NewDirsFactory(dirFactory DirFactory, countFactory func() int) DirsFactory {
	return func(parent *Dir, depth int) []*Dir {
		count := countFactory()
		dirs := make([]*Dir, count)
		for i := range dirs {
			dirs[i] = dirFactory(parent, depth)
		}
		return dirs
	}
}

func NullFilesFactory() FilesFactory {
	return func(*Dir) []*File {
		return nil
	}
}

func NullDirsFactory() DirsFactory {
	return func(*Dir, int) []*Dir {
		return nil
	}
}
