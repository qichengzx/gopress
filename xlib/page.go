package xlib

type PageNav struct {
	PageCount int
	PageSlice []int
}

// TODO need a better way to handle pagnition
func (pn *PageNav) Handler() *PageNav {
	page := make([]int, pn.PageCount)
	for i := 0; i < pn.PageCount; i++ {
		page[i] = i + 1
	}
	pn.PageSlice = page

	return pn
}
