package filetest

import (
	"testing"

	"github.com/ddirect/xrand"
)

type MinMax struct {
	Min int
	Max int
}

type TreeOptions struct {
	CharSet   string
	CharCount MinMax
	FileCount MinMax
	DirCount  MinMax
	Depth     int
}

func DefaultTreeOptions() TreeOptions {
	return TreeOptions{
		LowerCaseChars,
		MinMax{16, 32},
		MinMax{2, 4},
		MinMax{2, 4},
		3,
	}
}

func FlatTreeOptions(minFiles, maxFiles int) TreeOptions {
	return TreeOptions{
		LowerCaseChars,
		MinMax{16, 32},
		MinMax{minFiles, maxFiles},
		MinMax{},
		0,
	}
}

func NewRandomTree(rnd *xrand.Xrand, o TreeOptions) *Dir {
	res, _ := NewRandomTree2(rnd, o)
	return res
}

// returns also the name factory
func NewRandomTree2(rnd *xrand.Xrand, o TreeOptions) (*Dir, func() string) {
	nameFactory := NewRandomNameFactory(rnd, rnd.UniformFactory(o.CharCount.Min, o.CharCount.Max), o.CharSet)
	entryFactory := NewEntryFactory(nameFactory)
	fileFactory := NewFileFactory(entryFactory)
	filesFactory := NewFilesFactory(fileFactory, rnd.UniformFactory(o.FileCount.Min, o.FileCount.Max))

	dirsFactory, dfSet := NewDirsFactory(o.Depth, rnd.UniformFactory(o.DirCount.Min, o.DirCount.Max))
	dfSet(NewDirFactory(entryFactory, filesFactory, dirsFactory))
	return NewDirFactory(NullEntryFactory(), filesFactory, dirsFactory)(nil, 0), nameFactory
}

func CommitNewRandomTree(t *testing.T, altRoot string, o TreeOptions, m Mixes) (string, *Dir, DirStats) {
	rnd := xrand.New()
	dest := TempDir(t, altRoot)
	tree := NewRandomTree(rnd, o)
	return dest, tree, CommitMixed(rnd, tree, m, dest)
}

func CommitNewDefaultRandomTree(t *testing.T) (string, *Dir, DirStats) {
	return CommitNewRandomTree(t, "gotest_tree", DefaultTreeOptions(), DefaultMixes())
}
