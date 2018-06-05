package post

func SliceToMAP(sl []string) map[string]int {
	stmap := make(map[string]int)

	for _, str := range sl {
		index, ok := stmap[str]

		if ok {
			index++
			stmap[str] = index
		} else {
			stmap[str] = 1
		}

	}
	return stmap
}
