package assetshandler

import (
	_ "embed"
	"net/http"
)

//go:embed alpine3.min.js
var alpineJS string

//go:embed bootstrap5.2.3.min.css
var bootstrapCSS string

//go:embed favicon.ico
var favicon string

//go:embed bootstrap5.2.3.bundle.min.js
var boostrapBundleJS string

type response struct {
	body        string
	statusCode  int
	contentType string
}

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	resp := &response{
		statusCode:  http.StatusOK,
		contentType: "text/html",
	}
	switch r.RequestURI {
	case "/_assets/alpine.js":
		resp.contentType = "application/javascript"
		resp.body = alpineJS
	case "/_assets/bootstrap.css":
		resp.contentType = "text/css"
		resp.body = bootstrapCSS
	case "/_assets/bootstrap.bundle.js":
		resp.contentType = "application/javascript"
		resp.body = boostrapBundleJS
	case "/favicon.ico":
		resp.contentType = "image/x-icon"
		resp.body = favicon
	default:
		resp.statusCode = http.StatusNotFound
		resp.body = r.RequestURI + " not found\n"
	}
	w.Header().Set("Content-Type", resp.contentType)
	w.WriteHeader(resp.statusCode)
	if _, err := w.Write([]byte(resp.body)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
