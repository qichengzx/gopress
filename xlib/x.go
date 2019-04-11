package xlib

import (
	"bytes"
	"github.com/qichengzx/gopress/config"
	"github.com/qichengzx/gopress/generator"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Site struct {
	Posts      []generator.Post
	CatPosts   map[string][]generator.Post
	TagPosts   map[string][]generator.Post
	Archive    Archive
	Categories map[string]int
	Tags       map[string]int
	Recent     []generator.Post

	CurrentPage      string
	CurrentPageTitle string
	CurrentPageIndex int
	PrevPageIndex    int
	NextPageIndex    int
	PageNav          *PageNav

	CurrentPost generator.Post

	Cfg *config.Config

	CopyRight string
}

type Archive struct {
	Year     map[string][]generator.Post
	Archives map[string][]generator.Post
}

const (
	indexPage = "index.html"

	PageTypeIndex  = "index"
	PageTypeTag    = "tag"
	PageTypeCat    = "category"
	PageTypePost   = "post"
	PageTypeArh    = "archives"
	PageTypeArhIdx = "archiveIndex"
)

func New(cfFile string) *Site {
	var cfg = config.NewProvider(cfFile)

	appPath, _ := os.Getwd()
	postPath := filepath.Join(appPath, cfg.SourceDir)

	pw, tags, cates := generator.GetPosts(postPath, cfg)
	var Recent []generator.Post
	if len(pw.Posts) > 5 {
		Recent = pw.Posts[:5]
	} else {
		Recent = pw.Posts
	}

	return &Site{
		Posts:    pw.Posts,
		CatPosts: pw.CatPosts,
		TagPosts: pw.TagPosts,
		Archive: Archive{
			Archives: pw.Archives,
		},
		Categories: generator.SliceToMAP(cates),
		Tags:       generator.SliceToMAP(tags),
		Recent:     Recent,

		CurrentPage:      PageTypeIndex,
		CurrentPageIndex: 1,

		Cfg:       cfg,
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
	generator.ClearDir(s.Cfg.PublicDir)

	bt := s.renderPage()
	generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, indexPage))

	if s.PageNav.PageCount > 0 {
		for i := 1; i <= s.PageNav.PageCount; i++ {
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
			generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.PaginationDir, p, indexPage))
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
				p.PostNav = generator.PostNav{Next: generator.Nav{}, Prev: generator.Nav{}}
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

		generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, p.Path))
	}

	//TODO 分类，标签 暂不处理分页

	//分类页
	s.CurrentPage = PageTypeCat
	for cat, posts := range s.CatPosts {
		s.Posts = posts
		s.CurrentPageTitle = cat

		bt = s.renderPage()
		generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.CategoryDir, cat, indexPage))
	}

	//标签页
	s.CurrentPage = PageTypeTag
	for tag, posts := range s.TagPosts {
		s.Posts = posts
		s.CurrentPageTitle = tag

		bt = s.renderPage()
		generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.TagDir, tag, indexPage))
	}

	yearArchive := generator.GenArchive(posts)
	//Archived by year
	s.CurrentPage = PageTypeArh
	for year, posts := range yearArchive {
		s.Posts = posts
		s.CurrentPageTitle = year

		bt = s.renderPage()
		generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.ArchiveDir, year, indexPage))
	}

	//Archived by month
	s.CurrentPage = PageTypeArh
	for m, posts := range s.Archive.Archives {
		s.Posts = posts
		s.CurrentPageTitle = m

		bt = s.renderPage()
		generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.ArchiveDir, m, indexPage))
	}

	//Archive Index Page
	s.CurrentPage = PageTypeArhIdx
	s.Archive.Year = yearArchive
	bt = s.renderPage()
	generator.WriteFile(bt, filepath.Join(s.Cfg.PublicDir, s.Cfg.ArchiveDir, indexPage))

	s.copyAsset()
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

	var t = filepath.Join(config.ThemeDir, s.Cfg.Theme, "/layout/*.html")
	tmpl, err := template.ParseGlob(t)
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(&doc, "layout", s)

	return doc.Bytes()
}

func (s Site) copyAsset() {
	assets := map[string]string{
		filepath.Join(s.Cfg.SourceDir, "images"):                       "images",
		filepath.Join(config.ThemeDir, s.Cfg.Theme, "source/css"):      "css",
		filepath.Join(config.ThemeDir, s.Cfg.Theme, "source/js"):       "js",
		filepath.Join(config.ThemeDir, s.Cfg.Theme, "source/fancybox"): "fancybox",
	}

	for src, dst := range assets {
		generator.CopyDir(src, filepath.Join(s.Cfg.PublicDir, dst))
	}
}

func copyRight() string {
	return time.Now().Format("2006")
}
