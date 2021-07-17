package filetest

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ddirect/check"
	"github.com/ddirect/xrand"
)

var validChars = makeValidChars()

type FileFactory func(basePath string, dirPath string)
type DirFactory func(basePath string, dirPath string, depth int)

func makeValidChars() []byte {
	var valid []byte
	addRange := func(min, max byte) {
		for i := min; i <= max; i++ {
			valid = append(valid, i)
		}
	}
	add := func(s string) {
		for i := 0; i < len(s); i++ {
			valid = append(valid, s[i])
		}
	}
	addRange('a', 'z')
	addRange('A', 'Z')
	addRange('0', '9')
	add(" ()~#_-")
	return valid
}

func NewRandomNameFactory(rnd xrand.Xrand, sizeFactory func() int) func() string {
	return func() string {
		var b strings.Builder
		size := sizeFactory()
		for i := 0; i < size; i++ {
			b.WriteByte(validChars[rnd.Intn(len(validChars))])
		}
		return b.String()
	}
}

func NewRandomFileFactory(rnd xrand.Xrand, nameFactory func() string, sizeFactory func() int, newFileCallback func(string)) FileFactory {
	buf := make([]byte, 0x10000)
	slice := func(size int) []byte {
		if size < len(buf) {
			return buf[:size]
		}
		return buf
	}
	return func(base, dirPath string) {
		newPath := filepath.Join(dirPath, nameFactory())
		file, err := os.OpenFile(filepath.Join(base, newPath), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
		check.E(err)
		defer file.Close() // fail safe
		todo := sizeFactory()
		for todo > 0 {
			todo -= check.IE(file.Write(rnd.Fill(slice(todo))))
		}
		f1 := file
		file = nil
		check.E(f1.Close())
		newFileCallback(newPath)
	}
}

func NewDirFactory(maxDepth int, nameFactory func() string, fileFactory FileFactory, dirFactory DirFactory, newDirCallback func(string)) DirFactory {
	return func(base string, path string, depth int) {
		newPath := filepath.Join(path, nameFactory())
		check.E(os.Mkdir(filepath.Join(base, newPath), 0775))
		fileFactory(base, newPath)
		if depth < maxDepth {
			dirFactory(base, newPath, depth+1)
		}
		newDirCallback(newPath)
	}
}

func NewMultipleFileFactory(fileFactory FileFactory, countFactory func() int) FileFactory {
	return func(base, dirPath string) {
		count := countFactory()
		for i := 0; i < count; i++ {
			fileFactory(base, dirPath)
		}
	}
}

func NewMultipleDirFactory(dirFactory DirFactory, countFactory func() int) DirFactory {
	return func(base string, dirPath string, depth int) {
		count := countFactory()
		for i := 0; i < count; i++ {
			dirFactory(base, dirPath, depth)
		}
	}
}
