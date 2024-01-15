package request

import "net/http"

type MetaData struct {
	Encoding   string
	RemoteAddr string
	UserAgent  string
}

func BuildMetaData(r *http.Request) (meta MetaData) {
	meta = MetaData{
		Encoding:   r.Header.Get("Content-Encoding"),
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
	}
	return meta
}
