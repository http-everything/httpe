package response

import (
	"html/template"
	"strings"
)

var tplFunc = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
}
