package filetest

import "github.com/ddirect/xrand"

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

func NewRandomTree(rnd *xrand.Xrand, o TreeOptions) *Dir {
	entryFactory := NewEntryFactory(NewRandomNameFactory(rnd, rnd.UniformFactory(o.CharCount.Min, o.CharCount.Max), o.CharSet))
	fileFactory := NewFileFactory(entryFactory)
	filesFactory := NewFilesFactory(fileFactory, rnd.UniformFactory(o.FileCount.Min, o.FileCount.Max))

	dirsFactory, dfSet := NewDirsFactory(o.Depth, rnd.UniformFactory(o.DirCount.Min, o.DirCount.Max))
	dfSet(NewDirFactory(entryFactory, filesFactory, dirsFactory))
	return NewDirFactory(NullEntryFactory(), filesFactory, dirsFactory)(nil, 0)
}
