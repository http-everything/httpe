package servedirectory

import "net/http"

func Handle(path string, dir string) http.Handler {
	return http.StripPrefix(path, http.FileServer(http.Dir(dir)))
}
