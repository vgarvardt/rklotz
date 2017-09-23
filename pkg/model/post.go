package model

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

const (
	postMetaSeparator = "+++"
)

var (
	ErrorUnknownFormat    = errors.New("Unknown post format")
	ErrorBadPostStructure = errors.New("Bad post structure: must be post meta lines, separator, post body. Separator: " + postMetaSeparator)
	ErrorBadMetaStructure = errors.New("Bad post meta structure, must have the following lines: post title, publishing date, post tags")
)

type formatHandler func(input string) string

var formatsMap map[string]formatHandler = map[string]formatHandler{
	"md": func(input string) string {
		html := string(blackfriday.Run([]byte(input)))
		// open all links in new tab
		html = strings.Replace(html, `<a href=`, `<a target="_blank" href=`, -1)
		// fix code class to make highlight.js work
		re := regexp.MustCompile(`<code class="language-(\w+)">`)
		html = re.ReplaceAllString(html, "<code class=\"$1\">")

		return html
	},
}

type Post struct {
	Path        string `storm:"id"`
	ID          string `storm:"unique"`
	Title       string
	PublishedAt time.Time `storm:"index"`
	Tags        []string
	Body        string
	Format      string
	HTML        string
}

func NewPostFromFile(basePath, postPath string) (*Post, error) {
	post := &Post{Path: postPath[len(basePath) : len(postPath)-len(filepath.Ext(postPath))]}

	fileContents, err := ioutil.ReadFile(postPath)
	if err != nil {
		return nil, err
	}

	postParts := strings.SplitN(string(fileContents), postMetaSeparator, 2)
	if len(postParts) != 2 {
		return nil, ErrorBadPostStructure
	}

	postMeta := strings.Split(strings.TrimSpace(postParts[0]), "\n")
	if len(postMeta) < 3 {
		return nil, ErrorBadMetaStructure
	}

	post.Title = postMeta[0]

	post.PublishedAt, err = time.Parse(time.RFC822Z, postMeta[1])
	if nil != err {
		return nil, err
	}

	post.Tags = strings.Split(postMeta[2], ",")
	for i, tag := range post.Tags {
		post.Tags[i] = strings.TrimSpace(tag)
	}

	post.Body = strings.TrimSpace(postParts[1])
	post.Format = strings.ToLower(filepath.Ext(postPath)[1:])
	if handler, ok := formatsMap[post.Format]; ok {
		post.HTML = handler(post.Body)
	} else {
		return nil, ErrorUnknownFormat
	}

	h := sha1.New()
	h.Write([]byte(post.Path))
	post.ID = fmt.Sprintf("%x", h.Sum(nil))

	return post, nil
}
