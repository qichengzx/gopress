package sitemap

import (
	"bytes"
	"path/filepath"
	"text/template"
)

var indexSitemap = "sitemap.xml"
var indexSitemapTmpl = `<?xml version="1.0" encoding="UTF-8"?><?xml-stylesheet type="text/xsl" href="sitemap.xsl"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	{{ $domain := .SiteDomain }}
	{{ range $M := .IndexMap }}
    <sitemap>
        <loc>{{ $domain }}/{{ $M.Name }}</loc>
        <lastmod>{{ $M.Lastmod }}</lastmod>
    </sitemap>
	{{ end }}
</sitemapindex>`

//TODO Lastmod 
func (s Sitemap) indexSitemap() {
	var bf bytes.Buffer
	var indexMap = []IndexMap{}

	for _, name := range maps {
		var M = IndexMap{
			Name:    name,
			Lastmod: "",
		}

		indexMap = append(indexMap, M)
	}

	s.IndexMap = indexMap
	t := template.New("indexSitemap")
	t, _ = t.Parse(indexSitemapTmpl)
	t.Execute(&bf, s)

	makeFile(bf.Bytes(), filepath.Join(s.RootPath, indexSitemap))
}
