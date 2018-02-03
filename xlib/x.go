package xlib

import (
	"bytes"
	"html/template"
	"math"
	"path/filepath"
	"x/config"
	"x/post"
	"strconv"
)

type Site struct {
	Posts      []post.Post
	Categories map[string]int
	Tags       map[string]int

	CurrentPage      string
	CurrentPageIndex int
	PrevPageIndex    int
	NextPageIndex    int
	PageNav          *PageNav

	Cfg      *config.Config
	ThemeCfg *config.ThemeCfg
}

var indexPage = "index.html"

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
	s.CurrentPageIndex = 1

	// backup
	var posts = s.Posts
	if pn.PageCount > 1 {
		s.NextPageIndex = 2
		s.Posts = posts[:perPage]
	}

	clearDir(s.Cfg.PublicDir)
	bt := s.renderPage()

	makeFile(bt, filepath.Join(s.Cfg.PublicDir, indexPage))


	if pageCount > 0 {
		for i := perPage; i <= s.PageNav.PageCount; i++ {
			s.Posts = posts[i*perPage-perPage : i*perPage]
			s.CurrentPageIndex = i
			s.NextPageIndex = i + 1
			s.PrevPageIndex = i - 1
			bt := s.renderPage()

			p := strconv.Itoa(i)

			makeFile(bt,filepath.Join(s.Cfg.PublicDir,s.Cfg.PaginationDir, p, indexPage))
		}
	}
}

func (s *Site) renderPage() []byte {
	var doc bytes.Buffer

	var t = filepath.Join(s.Cfg.ThemeDir, s.Cfg.Theme, "/layout/*.html")
	tmpl, err := template.ParseGlob(t)
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(&doc, "layout", s)

	return []byte(doc.String())
}
