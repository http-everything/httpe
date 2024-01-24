Form Input:
{{- range $key,$value := .Input.Form }}
* {{ $key }}: {{ $value }}
{{- end }}

Json Input:
{{- range $key,$value := .Input.JSON }}
* {{ $key }}: {{ $value }}
{{- end }}

URL Parameters:
{{- range $key,$value := .Input.Params }}
* {{ $key }}: {{ $value }}
{{- end }}

Headers:
{{- range $key,$values := .Meta.Headers }}
* {{ $key }}: {{ $values }}
{{- end }}

First file upload:
{{ $upload := index .Input.Uploads 0 }}
Field Name: {{ $upload.FieldName }}
File Name: {{ $upload.FileName }}
Stored: {{ $upload.Stored }}
Type: {{ $upload.Type }}
Size: {{ $upload.Size }}

Meta Data:
* User Agent: {{ .Meta.UserAgent }}
* Method: {{ .Meta.Method }}
* URL: {{ .Meta.URL }}