package main

import (
	"time"
)

type Article struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Date string

func Today() Date {
	return Date(time.Now().Format("2006-01-02"))
}

type ReadingList struct {
	Articles map[Date][]Article `json:"articles"`
}

func (rl *ReadingList) AddArticle(date Date, article Article) {
	if rl.Articles == nil {
		rl.Articles = make(map[Date][]Article)
	}
	rl.Articles[date] = append(rl.Articles[date], article)
}
