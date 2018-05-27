package sitemap

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var maps = []string{"post-sitemap.xml", "page-sitemap.xml", "category-sitemap.xml", "tag-sitemap.xml"}

type Sitemap struct {
	RootPath   string
	SiteDomain string
	//Maps is for index-sitemap.xml

	//Items is for category,tag,etc...
	IndexMap []IndexMap
	Items    []Item
}

type IndexMap struct {
	Name    string
	Lastmod string
}

type Item struct {
	Permalink  string
	Lastmod    string
	Changefreq string
	Priority   float32
}

func NewRender(path, domain string) Sitemap {
	var Sitemap = Sitemap{
		RootPath:   path,
		SiteDomain: domain,
	}

	return Sitemap
}

func (sm Sitemap) Go() {
	sm.indexSitemap()
}

func makeFile(c []byte, file string) {
	dir := filepath.Dir(file)
	os.MkdirAll(dir, 0777)

	err := ioutil.WriteFile(file, c, 0644)
	if err != nil {
		panic(err)
	}
}
