package formatter

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var (
	// ErrorUnknownFormat is the error returned when trying to load a post with unknown format
	ErrorUnknownFormat = errors.New("unknown post format")
)

// Formatter defines type for formatter
type Formatter func(raw, format string) (string, error)

// New instantiates new Formatter
func New() Formatter {
	return newImpl().format
}

type implFormatter struct {
	mdHTML goldmark.Markdown
}

func newImpl() *implFormatter {
	return &implFormatter{
		mdHTML: goldmark.New(
			goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(),
		),
	}
}

func (f *implFormatter) format(raw, format string) (string, error) {
	if format == "md" {
		return f.formatMD(raw), nil
	}

	return "", ErrorUnknownFormat
}

var (
	reHref = regexp.MustCompile(`<a href="(.+[^"])">`)
	reCode = regexp.MustCompile(`<code class="language-(\w+)">`)
)

func (f *implFormatter) formatMD(raw string) string {
	var buf bytes.Buffer
	if err := f.mdHTML.Convert([]byte(raw), &buf); err != nil {
		return fmt.Sprintf("<p>Error rendering Markdown: <code>%s</code></p>", err.Error())
	}

	rendered := buf.String()
	// TODO: there must be a way to do this using plugins
	// add target-_blank to urls
	rendered = reHref.ReplaceAllString(rendered, `<a href="$1" target="_blank">`)
	// fix code class to make highlight.js work
	rendered = reCode.ReplaceAllString(rendered, `<code class="$1">`)

	return rendered
}
