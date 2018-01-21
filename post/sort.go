package post

import (
	"sort"
	"time"
)

type PostWrapper struct {
	post []Post
	by   func(p, q *Post) bool
}

type SortBy func(p, q *Post) bool

func (pw PostWrapper) Len() int {
	return len(pw.post)
}

func (pw PostWrapper) Swap(i, j int) {
	pw.post[i], pw.post[j] = pw.post[j], pw.post[i]
}

func (pw PostWrapper) Less(i, j int) bool {
	return pw.by(&pw.post[i], &pw.post[j])
}

func SortPost(posts []Post) []Post {
	sort.Sort(PostWrapper{posts, func(p, q *Post) bool {
		t1, _ := time.Parse("2006-01-02 15:04:05", p.Date)
		t2, _ := time.Parse("2006-01-02 15:04:05", q.Date)

		p.Unixtime = t1.Unix()
		q.Unixtime = t2.Unix()

		return q.Unixtime < p.Unixtime
	}})

	return posts
}
