package post

import (
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
}

type Tag struct {
	Name string
}

var (
	contentLine = 6
	fileExt     = ".md"
)

func GetPosts(path string) ([]Post, []string, []string) {
	return getPostlist(path)
}

func getPostlist(path string) ([]Post, []string, []string) {
	var Posts []Post
	var tags []string
	var cates []string

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

		p.setContent(path)
		p.setYear()
		p.setUnixtime()

		Posts = append(Posts, p)
		cates = append(cates, p.Category)

		for _, t := range p.Tags {
			tags = append(tags, t)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return SortPost(Posts), tags, cates
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

func (p *Post) setLink(l string, id string) *Post {
	now, _ := time.Parse("2006-01-02 15:04:05", p.Date)

	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	r := strings.NewReplacer(
		":year", year,
		":month", month,
		":day", day,
		":id", id,
		":title", id,
		":category", p.Category)

	p.Link = r.Replace(l)
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
