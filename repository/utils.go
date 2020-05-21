package repository

// isShareRootPath check whether path is shared with root.
func isShareRootPath(path, rootPath []int) bool {
	if len(path) < len(rootPath) {
		return false
	}

	return isSamePath(path[:len(rootPath)], rootPath)
}

// isSamePath checks whether two paths are same.
func isSamePath(path1, path2 []int) bool {
	if len(path1) != len(path2) {
		return false
	}

	for i := 0; i < len(path1); i++ {
		if path1[i] != path2[i] {
			return false
		}
	}

	return true
}

// mergePath merges two path.
func mergePath(path1, path2 []int) []int {
	newPath := []int{}
	newPath = append(newPath, path1...)
	newPath = append(newPath, path2...)

	return newPath
}
