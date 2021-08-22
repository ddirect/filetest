package filetest

import (
	"github.com/ddirect/format"
)

type DirStats struct {
	TotalFiles   int
	ClonedFiles  int
	LinkedFiles  int
	UniqueHashes int
}

func (ds DirStats) String() string {
	return ds.AppendToTable(DirStatsTable()).String()
}

func DirStatsTable() *format.Table {
	return new(format.Table).AppendColumn(
		"total files",
		"unique hashes",
		"cloned files",
		"linked files",
	)
}

func (ds *DirStats) AppendToTable(t *format.Table) *format.Table {
	t.AppendColumn(
		ds.TotalFiles,
		ds.UniqueHashes,
		ds.ClonedFiles,
		ds.LinkedFiles,
	)
	return t
}

func (ds *DirStats) Merge(o DirStats) {
	ds.TotalFiles += o.TotalFiles
	ds.ClonedFiles += o.ClonedFiles
	ds.LinkedFiles += o.LinkedFiles
	ds.UniqueHashes += o.UniqueHashes
}
