package renderer

import "net/http"

// Renderer is the interface for rendering data to a client in required format
type Renderer interface {
	// Render renders the data with response code to a HTTP response writer
	Render(w http.ResponseWriter, code int, data interface{})
}
