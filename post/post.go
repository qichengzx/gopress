package post

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/qichengzx/gopress/config"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type PostWarp struct {
	Posts    []Post
	CatPosts map[string][]Post
	TagPosts map[string][]Post
	Archives map[string][]Post
}

type Post struct {
	ID        string
	Title     string
	Category  string
	Created   time.Time
	Date      string
	Year      int
	Unixtime  int64
	Tags      []string
	Content   template.HTML
	Permalink string
	Path      string
	Index     int

	PostNav PostNav
}

type PostNav struct {
	Prev Nav
	Next Nav
}

type Nav struct {
	Title string
	Link  string
}

type Tag struct {
	Name string
}

var (
	contentLine = 6
	fileExt     = ".md"

	myCfg *config.Config
)

const (
	postDir  = "_posts"
	draftDir = "_draft"
)

func GetPosts(path string, cfg *config.Config) (PostWarp, []string, []string) {
	myCfg = cfg
	return getPostlist(path)
}

func getPostlist(path string) (PostWarp, []string, []string) {
	var (
		pw    PostWarp
		tags  []string
		cates []string
	)
	var cat = map[string][]Post{}
	var tag = map[string][]Post{}
	var arh = map[string][]Post{}

	path = filepath.Join(path, postDir)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		var p = Post{}

		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != fileExt {
			return nil
		}

		fileID := fileName(f.Name())

		p.setID().
			setContent(path).
			setCreated().
			setYear().
			setUnixtime().
			setLink(fileID)

		if p.Category == "" {
			p.Category = myCfg.DefaultCategory
		}

		cat[p.Category] = append(cat[p.Category], p)

		m := formatMonth(p.Date)
		arh[m] = append(arh[m], p)

		pw.Posts = append(pw.Posts, p)
		cates = append(cates, p.Category)

		for _, t := range p.Tags {
			tags = append(tags, t)
			tag[t] = append(tag[t], p)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	pw.CatPosts = cat
	pw.Posts = SortPost(pw.Posts)
	pw.TagPosts = tag
	pw.Archives = arh
	return pw, tags, cates
}

func (p *Post) setContent(fileName string) *Post {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(content, &p)
	if err != nil {
		panic(err)
	}

	p.Content = getContent(content)

	return p
}

func (p *Post) setID() *Post {
	id := md5.New()
	id.Write([]byte(p.Title))

	p.ID = hex.EncodeToString(id.Sum(nil))

	return p
}

func (p *Post) setLink(fileName string) *Post {
	if p.Permalink != "" {
		p.Path = p.Permalink
		if myCfg.RelativeLink {
			p.Permalink = myCfg.Root + p.Permalink
		} else {
			p.Permalink = myCfg.URL + myCfg.Root + p.Permalink
		}

		return p
	}
	now, _ := time.Parse("2006-01-02 15:04:05", p.Date)

	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	r := strings.NewReplacer(
		":year", year,
		":month", month,
		":day", day,
		":id", p.ID,
		":title", fileName,
		":category", p.Category)

	p.Path = r.Replace(myCfg.Permalink)
	if myCfg.RelativeLink {
		p.Permalink = myCfg.Root + p.Path
	} else {
		p.Permalink = myCfg.URL + myCfg.Root + p.Path
	}

	return p
}

func (p *Post) setUnixtime() *Post {
	p.Unixtime = formatUnix(p.Date)
	return p
}

func (p *Post) setYear() *Post {
	p.Year = formatYear(p.Date)
	return p
}

func (p *Post) setCreated() *Post {
	p.Created = formatUTC(p.Date)
	return p
}

func (p *Post) SetNav(p1, p2 *Post) *Post {
	if p1 == nil {
		p.PostNav = PostNav{
			Prev: Nav{
				Title: p2.Title,
				Link:  p2.Permalink,
			},
			Next: Nav{},
		}
	} else if p2 == nil {
		p.PostNav = PostNav{
			Prev: Nav{},
			Next: Nav{
				Title: p1.Title,
				Link:  p1.Permalink,
			},
		}
	} else {

		p.PostNav = PostNav{
			Prev: Nav{
				Title: p2.Title,
				Link:  p2.Permalink,
			},
			Next: Nav{
				Title: p1.Title,
				Link:  p1.Permalink,
			},
		}
	}

	return p
}

func fileName(f string) string {
	r := []rune(f)
	length := len(r)
	return string(r[0 : length-3])
}

func getContent(c []byte) template.HTML {
	lines := strings.Split(string(c), "\n")
	content := strings.Join(lines[contentLine:], "\n")
	str := blackfriday.MarkdownCommon([]byte(content))

	return template.HTML(str)
}

func GenArchive(posts []Post) map[string][]Post {
	var arh = map[string][]Post{}

	for _, post := range posts {
		key := strconv.Itoa(post.Year)
		arh[key] = append(arh[key], post)
	}

	return arh
}

func formatMonth(layout string) string {
	t, err := time.Parse("2006-01-02 15:04:05", layout)
	if err != nil {
		panic(err)
	}
	return t.Format("2006/01")
}

func formatDate(layout string) string {
	t, err := time.Parse("2006-01-02 15:04:05", layout)
	if err != nil {
		panic(err)
	}
	return t.Format("2006-01-02")
}

func formatDatetime(layout string) string {
	t, err := time.Parse("2006-01-02 15:04:05", layout)
	if err != nil {
		panic(err)
	}
	return t.Format(time.RFC3339)
}

func formatYear(layout string) int {
	t, err := time.Parse("2006-01-02 15:04:05", layout)
	if err != nil {
		panic(err)
	}
	return t.Year()
}

func formatUnix(layout string) int64 {
	t, err := time.Parse("2006-01-02 15:04:05", layout)
	if err != nil {
		panic(err)
	}
	return t.Unix()
}

func formatUTC(layout string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", layout)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}
