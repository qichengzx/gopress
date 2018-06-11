package sitemap

import (
	"io/ioutil"
)

type Sitemap struct {
	RootPath   string
	SiteDomain string

	//Maps is for index-sitemap.xml
	IndexMap []IndexMap

	//Items is for category,tag,etc...
	Items []Item
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

func (sm Sitemap) Go(post, category, tag []Item) {
	sm.indexSitemap()
	sm.postSitemap(post)
	sm.categorySitemap(category)
	sm.tagSitemap(tag)
}

func makeFile(c []byte, file string) {
	err := ioutil.WriteFile(file, c, 0644)
	if err != nil {
		panic(err)
	}
}
