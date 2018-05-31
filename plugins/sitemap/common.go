package sitemap

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

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
	Permalink string
	Lastmod   string
}

func NewRender(path, domain string) Sitemap {
	var Sitemap = Sitemap{
		RootPath:   path,
		SiteDomain: domain,
	}

	return Sitemap
}

func (sm Sitemap) Go(post []Item, category []Item) {
	sm.indexSitemap()
	sm.postSitemap(post)
	sm.categorySitemap(category)
}

func makeFile(c []byte, file string) {
	dir := filepath.Dir(file)
	os.MkdirAll(dir, 0777)

	err := ioutil.WriteFile(file, c, 0644)
	if err != nil {
		panic(err)
	}
}
