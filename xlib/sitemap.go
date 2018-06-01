package xlib

import (
	"github.com/qichengzx/gopress/plugins/sitemap"
)

func (s *Site) postMap() []sitemap.Item {
	var items = make([]sitemap.Item, len(s.Posts))

	for _, post := range s.Posts {
		var item = sitemap.Item{
			Permalink: post.Link,
			Lastmod:   post.Date,
		}

		items = append(items, item)
	}

	return items
}

func (s *Site) categoryMap() []sitemap.Item {
	var items = make([]sitemap.Item, len(s.CatPosts))

	for cate, post := range s.CatPosts {
		var item = sitemap.Item{
			Permalink: cate,
			Lastmod:   post[len(post)-1].Date,
		}

		items = append(items, item)
	}

	return items
}
