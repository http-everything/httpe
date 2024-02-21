package requestdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"http-everything/httpe/pkg/rules"

	"http-everything/httpe/pkg/filetype"

	"github.com/gorilla/mux"

	"github.com/lithammer/shortuuid/v4"
)

const UploadPrefix = "httpe_upload_"

type Data struct {
	Meta  MetaData
	Input Input
}

type MetaData struct {
	RemoteAddr string
	UserAgent  string
	Method     string
	URL        string
	Headers    map[string]string
}

type Input struct {
	Form            Form
	JSON            JSON
	Params          Params
	Uploads         []Upload
	URLPlaceholders URLPlaceholders
}

type Upload struct {
	FieldName string
	FileName  string
	Size      int64
	Type      string
	Stored    string
}

type Form map[string]string
type JSON interface{}
type Params map[string]string
type URLPlaceholders map[string]string

func Collect(r *http.Request, ruleArgs rules.Args) (d Data, err error) {
	formInput := make(map[string]string)
	params := make(map[string]string)
	uploads := make([]Upload, 0)
	urlPlaceholders := make(map[string]string)
	var jsonInput interface{}
	d = Data{
		Meta: buildMetaData(r),
		Input: Input{
			Form:            formInput,
			JSON:            jsonInput,
			Params:          params,
			Uploads:         uploads,
			URLPlaceholders: urlPlaceholders,
		},
	}
	d.Input.Params, err = extractURLParams(r)
	if err != nil {
		return d, err
	}
	if r.Header.Get("Content-Length") == "0" {
		return d, nil
	}

	cType := r.Header.Get("Content-Type")
	// Extract json data
	if strings.HasPrefix(cType, "application/json") {
		d.Input.JSON, err = extractJSONInput(r)
		if err != nil {
			return d, err
		}
	}

	// Extract form data
	if strings.HasPrefix(cType, "multipart/form-data") {
		if d.Input.Form, err = extractMultipartFormData(r); err != nil {
			return d, err
		}
		if ruleArgs.FileUploads {
			if d.Input.Uploads, err = extractFileUploads(r); err != nil {
				return d, err
			}
		}
	} else {
		// Extract from Content-Type: application/x-www-form-urlencoded
		if d.Input.Form, err = extractFormInput(r); err != nil {
			return d, err
		}
	}

	// Extract placeholders from URL
	d.Input.URLPlaceholders = extractURLPlaceholders(r)
	return d, nil
}

func buildMetaData(r *http.Request) (meta MetaData) {
	URL, err := url.PathUnescape(r.URL.RequestURI())
	if err != nil {
		URL = r.URL.RequestURI()
	}
	meta = MetaData{
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Method:     r.Method,
		URL:        URL,
		Headers:    extractHeaders(r),
	}

	return meta
}

func extractHeaders(r *http.Request) (headers map[string]string) {
	headers = make(map[string]string)
	for name, values := range r.Header {
		for _, value := range values {
			headers[name] = value
		}
	}

	return headers
}

func extractURLParams(r *http.Request) (params Params, err error) {
	params = make(map[string]string)
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return params, fmt.Errorf("error parsing url: %w", err)
	}
	p, _ := url.ParseQuery(u.RawQuery)
	if err != nil {
		return params, fmt.Errorf("error parsing url parameters: %w", err)
	}
	for k, v := range p {
		params[k] = v[0]
	}

	return params, nil
}

func extractJSONInput(r *http.Request) (i interface{}, err error) {
	err = json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		return i, fmt.Errorf("error parsing json-data: %w", err)
	}

	return i, nil
}

func extractFormInput(r *http.Request) (i map[string]string, err error) {
	i = make(map[string]string)
	err = r.ParseForm()
	if err != nil {
		return i, fmt.Errorf("error parsing form-data: %w", err)
	}

	for k, v := range r.PostForm {
		for _, value := range v {
			i[k] = value
		}
	}

	return i, nil
}

func extractMultipartFormData(r *http.Request) (i map[string]string, err error) {
	i = make(map[string]string)
	err = r.ParseMultipartForm(0)
	if err != nil {
		return i, fmt.Errorf("error parsing multipart/form-data: %w", err)
	}

	for k, v := range r.PostForm {
		for _, value := range v {
			i[k] = value
		}
	}

	return i, nil
}

func extractFileUploads(r *http.Request) (uploads []Upload, err error) {
	uploads = make([]Upload, 0)
	if err = r.ParseMultipartForm(r.ContentLength); err != nil {
		return uploads, fmt.Errorf("error parsing form to extract uploads: %w", err)
	}

	m := r.MultipartForm
	for name, file := range m.File {
		upload, err := file[0].Open()
		if err != nil {
			return uploads, fmt.Errorf("error extracting uploads: %w", err)
		}
		defer upload.Close()

		//create a temporary file to store the uploaded file
		fn := os.TempDir() + "/" + UploadPrefix + shortuuid.New()
		dst, err := os.Create(fn)
		if err != nil {
			return uploads, fmt.Errorf("error creating temp file for upload: %w", err)
		}
		defer dst.Close()
		// copy the uploaded file to the destination
		if _, err := io.Copy(dst, upload); err != nil {
			return uploads, fmt.Errorf("error copying upload to destination: %w", err)
		}
		dst.Close()

		fType, err := filetype.Type(fn)
		if err != nil {
			return uploads, fmt.Errorf("error getting file type: %w", err)
		}

		// add the file to the list of uploaded files
		uploads = append(uploads, Upload{
			FieldName: name,
			FileName:  file[0].Filename,
			Size:      file[0].Size,
			Type:      fType,
			Stored:    fn,
		})
	}
	return uploads, nil
}

func extractURLPlaceholders(r *http.Request) (placeholders URLPlaceholders) {
	return mux.Vars(r)
}

func Mock() (d Data, err error) {
	var jd interface{}
	err = json.Unmarshal([]byte(`
		{
            "jkey1":"json value 1",
            "jkey2":"json value 2",
            "nested": {"nkey1":"nvalue1"}
        }`), &jd)
	if err != nil {
		return Data{}, err
	}
	return Data{
		Meta: MetaData{
			UserAgent:  "golang",
			URL:        "/some/path",
			RemoteAddr: "127.0.0.1",
			Method:     "get",
			Headers: map[string]string{
				"X-My-Header": "gotest",
				"Upper":       "upper",
				"lower":       "lower",
			},
		},
		Input: Input{
			Form: map[string]string{
				"Field1": "Field Value 1",
				"Field2": "Field Value 2",
			},
			Params: map[string]string{
				"Param1": "Param Value 1",
				"Param2": "Param Value 2",
			},
			URLPlaceholders: map[string]string{
				"URLVar1": "URL Value 1",
				"URLVar2": "URL Value 2",
				"redir":   "https://example.com",
			},
			JSON: jd,
			Uploads: []Upload{
				{
					FieldName: "my-upload-1",
					FileName:  "my-upload-1.txt",
					Size:      1,
					Type:      "text",
					Stored:    "/tmp/Coogh8cheiToowabili",
				},
			},
		},
	}, nil
}
