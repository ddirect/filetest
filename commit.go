package filetest

import (
	"github.com/ddirect/xrand"
)

type Zones struct {
	// in % - what remains are not committed
	NoChange, Shuffle, Reseed int
}

func DefaultZones() Zones {
	return Zones{25, 50, 75}
}

type Mixes struct {
	MaxSize int
	// in % - what remains are unique files
	Created, Cloned, Linked int
}

func DefaultMixes() Mixes {
	return Mixes{5000, 25, 50, 75}
}

func ShuffleFiles(rnd *xrand.Xrand, f []*File) {
	rnd.Shuffle(len(f), func(i, j int) {
		f[i], f[j] = f[j], f[i]
	})
}

func CommitDirs(tree *Dir, root string) int {
	return tree.EachDirRecursive(NewDirMakerFactory(root))
}

func CommitMixed(rnd *xrand.Xrand, tree *Dir, m Mixes, root string) DirStats {
	files := tree.AllFilesSlice()
	ShuffleFiles(rnd, files)
	CommitDirs(tree, root)
	return CommitFilesMixed(rnd, files, m, root)
}

// Commits files from a tree
// returns the files which have not been committed (excluded)
func CommitZonedFilesMixed(rnd1 *xrand.Xrand, rnd2 *xrand.Xrand, files []*File, z Zones, m Mixes, root string, stage2 bool) (DirStats, []*File) {
	noChangeLimit := z.NoChange * len(files) / 100
	shuffleLimit := z.Shuffle * len(files) / 100
	reseedLimit := z.Reseed * len(files) / 100
	// the rest is uncommitted
	if !stage2 {
		ShuffleFiles(rnd2, files)
	}
	ds := CommitFilesMixed(rnd1, files[:noChangeLimit], m, root)
	if stage2 {
		ShuffleFiles(rnd2, files[noChangeLimit:])
	}
	ds.Merge(CommitFilesMixed(rnd1, files[noChangeLimit:shuffleLimit], m, root))
	ds.Merge(CommitFilesMixed(rnd2, files[shuffleLimit:reseedLimit], m, root))
	return ds, files[reseedLimit:]
}

func CommitFilesMixed(rnd *xrand.Xrand, files []*File, m Mixes, root string) DirStats {
	if len(files) == 0 {
		return DirStats{}
	}
	createLimit := m.Created * len(files) / 100
	cloneLimit := m.Cloned * len(files) / 100
	linkLimit := m.Linked * len(files) / 100
	uniqueLimit := len(files)

	if createLimit == 0 {
		createLimit = 1
	}

	create := NewRandomFileFactory(rnd, root, rnd.UniformFactory(0, m.MaxSize))
	clone := NewCloneFileOperation(root, root)
	link := NewLinkFileOperation(root, root)

	i := 0

	var ds DirStats
	for ; i < createLimit; i++ {
		create(files[i])
		ds.UniqueHashes++
	}
	for ; i < cloneLimit; i++ {
		clone(files[rnd.Intn(createLimit)], files[i])
		ds.ClonedFiles++
	}
	for ; i < linkLimit; i++ {
		link(files[rnd.Intn(cloneLimit)], files[i])
		ds.LinkedFiles++
	}
	for ; i < uniqueLimit; i++ {
		create(files[i])
		ds.UniqueHashes++
	}
	ds.TotalFiles = ds.UniqueHashes + ds.ClonedFiles + ds.LinkedFiles
	return ds
}
