package xlib

import (
	"github.com/gorilla/feeds"
	"log"
	"path/filepath"
	"time"
)

func (s *Site) Atom() {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       s.Cfg.Title,
		Subtitle:    s.Cfg.SubTitle,
		Link:        &feeds.Link{Href: s.Cfg.URL},
		Description: s.Cfg.Description,
		Author:      &feeds.Author{Name: s.Cfg.Author},
		Created:     now,
	}

	//TODO Description is too long
	var items []*feeds.Item
	for _, post := range s.Posts {
		var item = &feeds.Item{
			Title:       post.Title,
			Id:          post.ID,
			Link:        &feeds.Link{Href: post.Permalink},
			Description: string(post.Content),
			Content:     string(post.Content),
			Created:     post.Created,
			Updated:     now,
		}

		items = append(items, item)
	}

	feed.Items = items
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	writeFile([]byte(atom), filepath.Join(s.Cfg.PublicDir, s.Cfg.Rss))
}
