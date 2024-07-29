package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/anaskhan96/soup"
)

func getArticleTitleTitleTag(doc *soup.Root) (string, error) {
	tag := doc.Find("title")
	if tag.Error != nil {
		return "", fmt.Errorf("could not find article title: %v", tag.Error)
	}
	title := tag.Text()
	if title == "" {
		return "", fmt.Errorf("could not find article title")
	}
	return title, nil
}

func getArticleTitleMetaTag(doc *soup.Root) (string, error) {
	tag := doc.Find("meta", "property", "og:title")
	if tag.Error != nil {
		return "", fmt.Errorf("could not find article title: %v", tag.Error)
	}
	title := tag.Attrs()["content"]
	if title == "" {
		return "", fmt.Errorf("could not find article title")
	}
	return title, nil
}

func getArticleTitle(url string) (string, error) {
	soup.Header("User-Agent", "curl/0.0.0")
	html, err := soup.Get(url)
	if err != nil {
		return "", err
	}
	doc := soup.HTMLParse(html)
	if title, err := getArticleTitleTitleTag(&doc); err == nil {
		return title, nil
	}
	if title, err := getArticleTitleMetaTag(&doc); err == nil {
		return title, nil
	}
	return "", fmt.Errorf("could not find article title")
}

func loadReadingList(path string) (*ReadingList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read reading list: %v", err)
	}

	var rl ReadingList
	if err := json.Unmarshal(data, &rl); err != nil {
		return nil, err
	}

	return &rl, nil
}

func addToReadingList(path, url, title string) (*ReadingList, error) {
	rl, err := loadReadingList(path)
	if err != nil {
		return nil, fmt.Errorf("could not load reading list: %v", err)
	}

	rl.AddArticle(Today(), Article{Title: title, URL: url})

	jsonData, err := json.MarshalIndent(rl, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not marshal reading list: %v", err)
	}

	return rl, os.WriteFile(path, jsonData, 0644)
}

func mdFilename(jsonPath string) string {
	jsonBase := filepath.Base(jsonPath)
	name := jsonPath[:len(jsonBase)-len(filepath.Ext(jsonPath))] + ".md"
	return filepath.Join(filepath.Dir(jsonPath), name)
}

func add(args []string) {
	if len(args) != 4 {
		log.Fatalf("usage: %s add <json list> <url>", os.Args[0])
	}

	var (
		list = args[2]
		url  = args[3]
	)

	if filepath.Ext(list) != ".json" {
		log.Fatalf("reading list must be a json file")
	}

	title, err := getArticleTitle(url)
	if err != nil {
		log.Fatalf("could not get article title: %v", err)
	}

	rl, err := addToReadingList(list, url, title)
	if err != nil {
		log.Fatalf("could not add article to reading list: %v", err)
	}

	log.Printf("writing to %s", mdFilename(list))
	mdFile, err := os.Create(mdFilename(list))
	defer mdFile.Close()
	if err != nil {
		log.Fatalf("could not open markdown file: %v", err)
	}

	if err := markdown.Execute(mdFile, rl.Articles); err != nil {
		log.Fatalf("could not write to markdown file: %v", err)
	}
}

func generate(args []string) {
	if len(args) != 3 {
		log.Fatalf("usage: %s generate <json list>", os.Args[0])
	}

	rl, err := loadReadingList(args[2])
	if err != nil {
		log.Fatalf("could not load reading list: %v", err)
	}

	markdown.Execute(os.Stdout, rl.Articles)
}

func baseUsage() {
	log.Fatalf("usage: %s [add|generate]", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		baseUsage()
	}
	switch os.Args[1] {
	case "add":
		add(os.Args)
	case "generate":
		generate(os.Args)
	default:
		baseUsage()
	}
}
