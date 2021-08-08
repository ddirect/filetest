package filetest

import (
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"github.com/ddirect/check"
	"github.com/ddirect/xrand"
	"golang.org/x/crypto/blake2b"
)

type File struct {
	Entry
	Hash []byte
}

func (f *File) String() string {
	return fmt.Sprintf("%02x %s", f.Hash, f.Path())
}

const newFilePerm = 0664
const newFileFlags = os.O_WRONLY | os.O_CREATE | os.O_EXCL

func hashEngine() hash.Hash {
	hash, err := blake2b.New256(nil)
	check.E(err)
	return hash
}

func makeBuffer() []byte {
	return make([]byte, 0x10000)
}

func NewRandomFileFactory(rnd *xrand.Xrand, base string, sizeFactory func() int) func(*File) {
	buf := makeBuffer()
	slice := func(size int) []byte {
		if size < len(buf) {
			return buf[:size]
		}
		return buf
	}
	hash := hashEngine()
	return func(f *File) {
		file, err := os.OpenFile(filepath.Join(base, f.Path()), newFileFlags, newFilePerm)
		check.E(err)
		defer file.Close() // fail safe
		todo := sizeFactory()
		hash.Reset()
		for todo > 0 {
			b := slice(todo)
			rnd.Fill(b)
			hash.Write(b)
			todo -= check.IE(file.Write(b))
		}
		f.Hash = hash.Sum(nil)
		f1 := file
		file = nil
		check.E(f1.Close())
	}
}

func NewCloneFileOperation(sbase, dbase string) func(*File, *File) {
	buf := makeBuffer()
	return func(sfile, dfile *File) {
		sf, err := os.Open(filepath.Join(sbase, sfile.Path()))
		check.E(err)
		defer check.DeferredE(sf.Close)
		df, err := os.OpenFile(filepath.Join(dbase, dfile.Path()), newFileFlags, newFilePerm)
		check.E(err)
		defer check.DeferredE(df.Close)
		check.I64E(io.CopyBuffer(df, sf, buf))
		dfile.Hash = sfile.Hash
	}
}

func NewLinkFileOperation(sbase, dbase string) func(*File, *File) {
	return func(sfile, dfile *File) {
		check.E(os.Link(filepath.Join(sbase, sfile.Path()), filepath.Join(dbase, dfile.Path())))
		dfile.Hash = sfile.Hash
	}
}

func NewFileFromStorageFactory(base string) func(Entry) *File {
	buf := makeBuffer()
	hash := hashEngine()
	return func(entry Entry) *File {
		file, err := os.Open(filepath.Join(base, entry.Path()))
		check.E(err)
		defer check.DeferredE(file.Close)
		hash.Reset()
		check.I64E(io.CopyBuffer(hash, file, buf))
		return &File{entry, hash.Sum(nil)}
	}
}
