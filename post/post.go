package post

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PostWarp struct {
	Posts    []Post
	CatPosts map[string][]Post
	TagPosts map[string][]Post
}

type Post struct {
	ID       string
	Title    string
	Category string
	Date     string
	Year     int
	Unixtime int64
	Tags     []string
	Content  template.HTML
	Link     string
	Index    int
}

type Tag struct {
	Name string
}

var (
	contentLine = 6
	fileExt     = ".md"

	Root      string
	Permalink string
)

func GetPosts(path string) (PostWarp, []string, []string) {
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

		p.setID()
		p.setContent(path)
		p.setYear()
		p.setUnixtime()
		p.setLink(fileID)

		cat[p.Category] = append(cat[p.Category], p)

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

	p.Link = Root + r.Replace(Permalink)
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

func fileName(f string) string {
	r := []rune(f)
	length := len(r)
	return string(r[0 : length-3])
}

func getContent(c []byte) template.HTML {
	lines := strings.Split(string(c), "\n")
	content := strings.Join(lines[contentLine:len(lines)], "\n")
	str := blackfriday.MarkdownCommon([]byte(content))

	return template.HTML(str)
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
