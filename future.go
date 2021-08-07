package filetest

func FutureDirsFactory(maxdepth int, f *DirsFactory) DirsFactory {
	return func(parentPath string, depth int) []*Dir {
		if depth > maxdepth {
			return nil
		}
		return (*f)(parentPath, depth)
	}
}
