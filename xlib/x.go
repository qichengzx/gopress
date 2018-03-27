package xlib

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"time"
	"gopress/config"
	"gopress/post"
)

type Site struct {
	Posts      []post.Post
	CatPosts   map[string][]post.Post
	TagPosts   map[string][]post.Post
	Categories map[string]int
	Tags       map[string]int

	CurrentPage      string
	CurrentPageTitle string
	CurrentPageIndex int
	PrevPageIndex    int
	NextPageIndex    int
	PageNav          *PageNav

	Cfg      *config.Config
	ThemeCfg *config.ThemeCfg

	CopyRight string
}

const (
	indexPage = "index.html"
	ThemeDir  = "themes"

	PageTypeIndex = "index"
	PageTypeTag   = "tag"
	PageTypeCat   = "category"
	PageTypePost  = "post"
)

func (s *Site) Build() {
	s.CurrentPage = PageTypeIndex
	postCount := len(s.Posts)

	s.CurrentPageIndex = 1
	s.makePagnition(postCount, s.Cfg.PerPage)
	s.copyRight()

	// backup
	var posts = s.Posts
	if s.PageNav.PageCount > 1 {
		s.NextPageIndex = s.CurrentPageIndex + 1
		s.Posts = posts[:s.Cfg.PerPage]
	}

	// TODO clear public dir only when generate page was success
	clearDir(s.Cfg.PublicDir)
	bt := s.renderPage()

	makeFile(bt, filepath.Join(s.Cfg.PublicDir, indexPage))

	if s.PageNav.PageCount > 0 {
		for i := s.Cfg.PerPage; i <= s.PageNav.PageCount; i++ {
			//It's not good enough
			lastIndex := 0
			if i*s.Cfg.PerPage > postCount {
				lastIndex = postCount
			} else {
				lastIndex = i * s.Cfg.PerPage
			}

			s.Posts = posts[i*s.Cfg.PerPage-s.Cfg.PerPage : lastIndex]
			s.CurrentPageIndex = i
			s.NextPageIndex = i + 1
			s.PrevPageIndex = i - 1
			bt := s.renderPage()

			p := strconv.Itoa(i)
			makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.PaginationDir, p, indexPage))
		}
	}

	//文章页
	s.Posts = posts
	s.CurrentPage = PageTypePost
	for i, p := range s.Posts {
		p.Index = i
		s.CurrentPageIndex = i
		bt := s.renderPage()

		makeFile(bt, filepath.Join(s.Cfg.PublicDir, p.Link))
	}

	//TODO 分类，标签 暂不处理分页

	//分类页
	s.CurrentPage = PageTypeCat
	for cat, posts := range s.CatPosts {
		s.Posts = posts
		s.CurrentPageTitle = cat

		bt := s.renderPage()
		makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.CategoryDir, cat, indexPage))
	}

	//标签页
	s.CurrentPage = PageTypeTag
	for tag, posts := range s.TagPosts {
		s.Posts = posts
		s.CurrentPageTitle = tag

		bt := s.renderPage()
		makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.TagDir, tag, indexPage))
	}

	s.style()
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

// TODO minify style.css
func (s *Site) style() {
	stylePath := filepath.Join(ThemeDir, s.Cfg.Theme, "/style.css")
	data, err := ioutil.ReadFile(stylePath)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(s.Cfg.PublicDir+"/style.css", data, 0644)
}

func (s *Site) copyRight() *Site {
	s.CopyRight = time.Now().Format("2006")
	return s
}
