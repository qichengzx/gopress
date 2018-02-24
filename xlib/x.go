package xlib

import (
	"bytes"
	"html/template"
	"math"
	"path/filepath"
	"strconv"
	"x/config"
	"x/post"
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
var ThemeDir = "themes"

func (s *Site) Build() {
	s.CurrentPage = "index"
	count := len(s.Posts)

	s.CurrentPageIndex = 1
	s.makePagnition(count, s.Cfg.PerPage)

	// backup
	var posts = s.Posts
	if s.PageNav.PageCount > 1 {
		s.NextPageIndex = 2
		s.Posts = posts[:s.Cfg.PerPage]
	}

	clearDir(s.Cfg.PublicDir)
	bt := s.renderPage()

	makeFile(bt, filepath.Join(s.Cfg.PublicDir, indexPage))

	if s.PageNav.PageCount > 0 {
		for i := s.Cfg.PerPage; i <= s.PageNav.PageCount; i++ {
			s.Posts = posts[i*s.Cfg.PerPage-s.Cfg.PerPage : i*s.Cfg.PerPage]
			s.CurrentPageIndex = i
			s.NextPageIndex = i + 1
			s.PrevPageIndex = i - 1
			bt := s.renderPage()

			p := strconv.Itoa(i)

			makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.PaginationDir, p, indexPage))
		}
	}
}

func (s *Site) makePagnition(count int, perPage int) *Site {
	pageCount := float64(0)

	if count > perPage {
		pageCount = math.Ceil(float64(count) / float64(perPage))
	}

	var pn = PageNav{}
	pn.PageCount = int(pageCount)
	s.PageNav = pn.Handler()

	return s
}

func (s *Site) renderPage() []byte {
	var doc bytes.Buffer

	var t = filepath.Join(ThemeDir, s.Cfg.Theme, "/layout/*.html")
	tmpl, err := template.ParseGlob(t)
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(&doc, "layout", s)

	return []byte(doc.String())
}
