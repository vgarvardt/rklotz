package renderer

import "net/http"

const (
	charsetUTF8                   = "charset=utf-8"
	MIMEApplicationXML            = "application/xml"
	MIMEApplicationXMLCharsetUTF8 = MIMEApplicationXML + "; " + charsetUTF8
	HeaderContentType             = "Content-Type"
)

type xmlRenderer struct{}

func NewXmlRenderer() *xmlRenderer {
	return &xmlRenderer{}
}

func (r *xmlRenderer) Render(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	w.WriteHeader(code)
	w.Write([]byte(data.(string)))
}
