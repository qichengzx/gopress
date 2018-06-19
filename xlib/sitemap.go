// sitemap.go fenerate aitemap.xml for site
package xlib

import (
	"github.com/qichengzx/gopress/plugins/sitemap"
)

func (s *Site) postMap() []sitemap.Item {
	var items = []sitemap.Item{}

	for _, post := range s.Posts {
		var item = sitemap.Item{
			Permalink: post.Permalink,
			Lastmod:   post.Date,
		}

		items = append(items, item)
	}

	return items
}

func (s *Site) categoryMap() []sitemap.Item {
	var items = []sitemap.Item{}

	for cate, posts := range s.CatPosts {
		var item = sitemap.Item{
			Permalink: cate,
			Lastmod:   posts[len(posts)-1].Date,
		}

		items = append(items, item)
	}

	return items
}

func (s *Site) tagMap() []sitemap.Item {
	var items = []sitemap.Item{}

	for tag, posts := range s.TagPosts {
		var item = sitemap.Item{
			Permalink: tag,
			Lastmod:   posts[len(posts)-1].Date,
		}

		items = append(items, item)
	}

	return items
}
