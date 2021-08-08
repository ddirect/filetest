package filetest

import (
	"reflect"
	"testing"

	"github.com/ddirect/xrand"
)

func createRandomTree(rnd *xrand.Xrand) *Dir {
	entryFactory := NewEntryFactory(NewRandomNameFactory(rnd, rnd.UniformFactory(200, 250), ValidChars))
	fileFactory := NewFileFactory(entryFactory)
	filesFactory := NewFilesFactory(fileFactory, rnd.UniformFactory(2, 4))

	var dirsFactory DirsFactory
	dirFactory := NewDirFactory(entryFactory, filesFactory, FutureDirsFactory(3, &dirsFactory))
	dirsFactory = NewDirsFactory(dirFactory, rnd.UniformFactory(2, 4))
	return NewDirFactory(NullEntryFactory(), filesFactory, dirsFactory)(nil, 0)
}

func TestCreateAndReloadTree(t *testing.T) {
	rnd := xrand.New()
	root := t.TempDir()
	tree := createRandomTree(rnd)
	tree.Sort()
	tree.EachDirRecursive(NewDirMakerFactory(root))
	tree.EachFileRecursive(NewRandomFileFactory(rnd, root, rnd.UniformFactory(100, 200)))
	st := NewDirFromStorage(root)
	st.Sort()
	if !reflect.DeepEqual(tree, st) {
		t.Fail()
	}
}
