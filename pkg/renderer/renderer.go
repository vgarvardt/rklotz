package renderer

import "net/http"

type Renderer interface {
	Render(w http.ResponseWriter, code int, data interface{})
}
