package filetest

func FutureDirsFactory(maxdepth int, f *DirsFactory) DirsFactory {
	return func(parent *Dir, depth int) []*Dir {
		if depth > maxdepth {
			return nil
		}
		return (*f)(parent, depth)
	}
}
