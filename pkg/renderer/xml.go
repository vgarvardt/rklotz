package renderer

import "net/http"

const (
	charsetUTF8                   = "charset=utf-8"
	mimeApplicationXML            = "application/xml"
	mimeApplicationXMLCharsetUTF8 = mimeApplicationXML + "; " + charsetUTF8
	headerContentType             = "Content-Type"
)

// XMLRenderer implements Renderer for XML content
type XMLRenderer struct{}

// NewXMLRenderer creates new XMLRenderer instance
func NewXMLRenderer() *XMLRenderer {
	return &XMLRenderer{}
}

// Render renders the data with response code to a HTTP response writer
func (r *XMLRenderer) Render(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(headerContentType, mimeApplicationXMLCharsetUTF8)
	w.WriteHeader(code)
	w.Write([]byte(data.(string)))
}
