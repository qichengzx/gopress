package post

import (
	"sort"
)

type postWrapper []Post

func (pw postWrapper) Len() int           { return len(pw) }
func (pw postWrapper) Less(i, j int) bool { return pw[i].Unixtime > pw[j].Unixtime }
func (pw postWrapper) Swap(i, j int)      { pw[i], pw[j] = pw[j], pw[i] }

func SortPost(posts []Post) []Post {
	sort.Sort(postWrapper(posts))

	return posts
}
