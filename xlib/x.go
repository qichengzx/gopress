package xlib

import (
	"bytes"
	"github.com/qichengzx/gopress/config"
	"github.com/qichengzx/gopress/plugins/sitemap"
	"github.com/qichengzx/gopress/post"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Site struct {
	Posts      []post.Post
	CatPosts   map[string][]post.Post
	TagPosts   map[string][]post.Post
	Archives   map[string][]post.Post
	Categories map[string]int
	Tags       map[string]int
	Recent     []post.Post

	CurrentPage      string
	CurrentPageTitle string
	CurrentPageIndex int
	PrevPageIndex    int
	NextPageIndex    int
	PageNav          *PageNav

	CurrentPost post.Post

	Cfg      *config.Config
	ThemeCfg *config.ThemeCfg

	CopyRight string
}

const (
	indexPage = "index.html"
	ThemeDir  = "themes"

	PageTypeIndex  = "index"
	PageTypeTag    = "tag"
	PageTypeCat    = "category"
	PageTypePost   = "post"
	PageTypeArh    = "archives"
	PageTypeArhIdx = "archiveIndex"
)

func New(cfFile string) *Site {
	var cfg = config.NewProvider(cfFile)
	var themeCfg = config.ThemeCfgProvider(filepath.Join(ThemeDir, cfg.Theme, cfFile))

	appPath, _ := os.Getwd()
	postPath := filepath.Join(appPath, cfg.SourceDir)

	pw, tags, cates := post.GetPosts(postPath, cfg)
	var Recent []post.Post
	if len(pw.Posts) > 5 {
		Recent = pw.Posts[:5]
	} else {
		Recent = pw.Posts
	}

	return &Site{
		Posts:      pw.Posts,
		CatPosts:   pw.CatPosts,
		TagPosts:   pw.TagPosts,
		Archives:   pw.Archives,
		Categories: post.SliceToMAP(cates),
		Tags:       post.SliceToMAP(tags),
		Recent:     Recent,

		CurrentPage:      PageTypeIndex,
		CurrentPageIndex: 1,

		Cfg:       cfg,
		ThemeCfg:  themeCfg,
		CopyRight: copyRight(),
	}
}

func (s *Site) Build() {
	s.CurrentPage = PageTypeIndex
	postCount := len(s.Posts)

	s.CurrentPageIndex = 1
	s.makePagnition(postCount)

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
		if i == 0 {
			if postCount > 1 {
				p.SetNav(nil, &s.Posts[i+1])
			} else {
				p.PostNav = post.PostNav{Next: post.Nav{}, Prev: post.Nav{}}
			}

		} else if i == postCount-1 {
			p.SetNav(&s.Posts[i-1], nil)
		} else {
			p.SetNav(&s.Posts[i-1], &s.Posts[i+1])
		}

		p.Index = i
		s.CurrentPageIndex = i
		s.CurrentPost = p

		bt = s.renderPage()

		makeFile(bt, filepath.Join(s.Cfg.PublicDir, p.Link))
	}

	//TODO 分类，标签 暂不处理分页

	//分类页
	s.CurrentPage = PageTypeCat
	for cat, posts := range s.CatPosts {
		s.Posts = posts
		s.CurrentPageTitle = cat

		bt = s.renderPage()
		makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.CategoryDir, cat, indexPage))
	}

	//标签页
	s.CurrentPage = PageTypeTag
	for tag, posts := range s.TagPosts {
		s.Posts = posts
		s.CurrentPageTitle = tag

		bt = s.renderPage()
		makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.TagDir, tag, indexPage))
	}

	yearArchive := post.GenArchive(posts)
	//Archived by year
	s.CurrentPage = PageTypeArh
	for year, posts := range yearArchive {
		s.Posts = posts
		s.CurrentPageTitle = year

		bt = s.renderPage()
		makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.ArchiveDir, year, indexPage))
	}

	//Archived by month
	s.CurrentPage = PageTypeArh
	for m, posts := range s.Archives {
		s.Posts = posts
		s.CurrentPageTitle = m

		bt = s.renderPage()
		makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.ArchiveDir, m, indexPage))
	}

	//Archive Index Page
	s.CurrentPage = PageTypeArhIdx
	s.Archives = yearArchive
	bt = s.renderPage()
	makeFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.ArchiveDir, indexPage))

	s.copyAsset()

	s.Posts = posts
	render := sitemap.NewRender(s.Cfg.PublicDir, s.Cfg.Url)
	render.Go(s.postMap(), s.categoryMap())

	s.Atom()
}

func (s *Site) makePagnition(count int) *Site {
	pageCount := float64(0)

	if count > s.Cfg.PerPage {
		pageCount = math.Ceil(float64(count) / float64(s.Cfg.PerPage))
	}

	var pn = PageNav{
		PageCount: int(pageCount),
	}
	s.PageNav = pn.Handler()

	return s
}

func (s Site) renderPage() []byte {
	var doc bytes.Buffer

	var t = filepath.Join(ThemeDir, s.Cfg.Theme, "/layout/*.html")
	tmpl, err := template.ParseGlob(t)
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(&doc, "layout", s)

	return doc.Bytes()
}

func (s Site) copyAsset() {
	err := CopyDir(filepath.Join(s.Cfg.SourceDir, "../images"), filepath.Join(s.Cfg.PublicDir, "images"))
	if err != nil {
		panic(err)
	}

	err = CopyDir(filepath.Join(ThemeDir, s.Cfg.Theme, "css"), filepath.Join(s.Cfg.PublicDir, "css"))
	if err != nil {
		panic(err)
	}
}

func copyRight() string {
	return time.Now().Format("2006")
}
