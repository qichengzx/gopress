package xlib

import (
	"github.com/qichengzx/gopress/plugins/sitemap"
)

func (s *Site)postMap() []sitemap.Item{
	var items = []sitemap.Item{}

	for _, post := range s.Posts {
		var item = sitemap.Item{
			Permalink:post.Link,
			Lastmod: post.Date,
		}

		items = append(items, item)
	}

	return items
}