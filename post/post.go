package post

import (
	"html/template"
	"time"
)

type Post struct {
	ID         string
	Title      string
	Category   string
	Date       string
	Datetime   time.Time
	Tags       []string
	Content    template.HTML
	Link       string
}
