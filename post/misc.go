package post

import "strings"

func WordToMAP(s string) map[string]int {
	var stmap = make(map[string]int)
	strs := strings.Fields(s)

	for _, str := range strs {
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
