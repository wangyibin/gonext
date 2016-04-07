package gonext

type PathNames []string
// PathNames func
func ParsePathNames(path string) PathNames {
	var pnames PathNames = []string{} // Param names
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)
		} else if path[i] == '*' {
			pnames = append(pnames, "_*")
		}
	}
	return pnames
}

func (pnames PathNames) contains(key string) bool {
	for _, pname := range pnames {
		if pname == key {
			return true
		}
	}
	return false
}