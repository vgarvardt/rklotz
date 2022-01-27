package formatter

import (
	"errors"
	"regexp"

	"github.com/russross/blackfriday/v2"
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
	mdHTMLRenderer *blackfriday.HTMLRenderer
}

func newImpl() *implFormatter {
	impl := new(implFormatter)

	// copied from CommonHTMLFlags with UseXHTML removed and HrefTargetBlank added
	htmlFlags := blackfriday.Smartypants |
		blackfriday.SmartypantsFractions | blackfriday.SmartypantsDashes | blackfriday.SmartypantsLatexDashes |
		blackfriday.HrefTargetBlank

	impl.mdHTMLRenderer = blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: htmlFlags,
	})

	return impl
}

func (f *implFormatter) format(raw, format string) (string, error) {
	switch format {
	case "md":
		return f.formatMD(raw), nil
	}

	return "", ErrorUnknownFormat
}

func (f *implFormatter) formatMD(raw string) string {
	html := string(blackfriday.Run([]byte(raw), blackfriday.WithRenderer(f.mdHTMLRenderer)))
	// fix code class to make highlight.js work
	re := regexp.MustCompile(`<code class="language-(\w+)">`)
	html = re.ReplaceAllString(html, "<code class=\"$1\">")

	return html
}
