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
	Tags     []string
	Content  template.HTML
	Link     string
}

var (
	contentLine = 6
	fileExt     = ".md"
)

func GetPosts(path string) []*Post {
	return getPostlist(path)
}

func getPostlist(path string) []*Post {
	var Posts []*Post
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		var p = &Post{}
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

		Posts = append(Posts, p)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return Posts
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

func getFileName(t string) string {
	r := []rune(t)
	length := len(r)
	return string(r[0 : length-3])
}

func getContent(c []byte) template.HTML {
	lines := strings.Split(string(c), "\n")
	content := strings.Join(lines[contentLine:len(lines)], "\n")
	str := blackfriday.MarkdownCommon([]byte(content))

	return template.HTML(str)
}
