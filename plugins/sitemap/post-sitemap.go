package sitemap

import (
	"bytes"
	"path/filepath"
	"text/template"
)

var postSitemap = "post-sitemap.xml"
var postSitemapTmpl = `<?xml version="1.0" encoding="UTF-8"?><?xml-stylesheet type="text/xsl" href="sitemap.xsl"?>
<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    {{ $domain := .SiteDomain }}
	<url>
        <loc>{{ $domain }}</loc>
        <changefreq>daily</changefreq>
        <priority>1</priority>
    </url>
	{{ range $item := .Items }}
    <url>
        <loc>{{ $domain }}{{ $item.Permalink }}</loc>
        <lastmod>{{ $item.Lastmod }}</lastmod>
        <changefreq>weekly</changefreq>
        <priority>0.6</priority>
    </url>
	{{ end }}
</urlset>`


func (s Sitemap) postSitemap(p []Item) {
	s.Items = p
	var bf bytes.Buffer
	t := template.New("postSitemap")
	t, _ = t.Parse(postSitemapTmpl)
	t.Execute(&bf, s)

	makeFile(bf.Bytes(), filepath.Join(s.RootPath, postSitemap))
}