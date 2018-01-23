package xlib

import (
	"fmt"
	"math"
	"x/config"
	"x/post"
)

type Site struct {
	Posts      []post.Post
	Categories map[string]int
	Tags       map[string]int

	CurrentPage      string
	CurrentPageIndex int
	PageNav          *PageNav

	Cfg *config.Config
}

func (s *Site) Build() {
	s.CurrentPage = "index"
	count := len(s.Posts)
	perPage := s.Cfg.PerPage
	pageCount := float64(0)

	if count > perPage {
		pageCount = math.Ceil(float64(count) / float64(perPage))
	}

	var pn = PageNav{}
	pn.PageCount = int(pageCount)
	s.PageNav = pn.Handler()

	fmt.Println(s)
}
