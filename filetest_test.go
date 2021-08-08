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
