package requesthandler

import (
	"fmt"
	"net/http"
)

func Get(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("http call")
	fmt.Fprint(w, "foo\n")
}
