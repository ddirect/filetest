package filetest

import (
	"testing"

	"github.com/ddirect/xrand"
)

func TestCreateAndReloadTree(t *testing.T) {
	rnd := xrand.New()
	root := t.TempDir()
	tree := NewRandomTree(rnd, TreeOptions{
		ValidChars,
		MinMax{200, 250},
		MinMax{2, 4},
		MinMax{2, 4},
		3,
	})
	tree.EachDirRecursive(NewDirMakerFactory(root))
	tree.EachFileRecursive(NewRandomFileFactory(rnd, root, rnd.UniformFactory(100, 200)))
	if !tree.Compare(NewDirFromStorage(root)) {
		t.Fail()
	}
}

func TestCommitZonedFilesMixed(t *testing.T) {
	rnd := xrand.New()
	root := t.TempDir()

	tree := NewRandomTree(rnd, TreeOptions{
		ValidChars,
		MinMax{32, 40},
		MinMax{0, 5},
		MinMax{3, 4},
		3,
	})

	files := tree.AllFilesSlice()
	dirCount := CommitDirs(tree, root)
	ds, excluded := CommitZonedFilesMixed(rnd, rnd, files, DefaultZones(), DefaultMixes(), root, true)
	t.Logf("%d dirs - %s - %d excluded files (already excluded from total)", dirCount, ds, len(excluded))
	tree.RemoveFiles(excluded)
	if !tree.Compare(NewDirFromStorage(root)) {
		t.Fail()
	}
}
