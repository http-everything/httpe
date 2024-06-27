package requestdata_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/http-everything/httpe/pkg/rules"

	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/share/extract"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	UserAgent         = "Test User Agent"
	Name              = "John Doe"
	Email             = "john@example.com"
	HeaderContentType = "Content-Type"
)

func TestRequestDataGet(t *testing.T) {
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "http",
			Host:     "localhost",
			Path:     "/test",
			RawQuery: "var1=foo&var2=bar",
		},
		Header: http.Header{
			"User-Agent": []string{UserAgent},
			"X-Test":     []string{"Test Header"},
		},
		Host: "localhost",
	}
	reqData, err := requestdata.Collect(req, rules.Args{})
	require.NoError(t, err)

	assert.Equal(t, UserAgent, reqData.Meta.UserAgent)
	assert.Equal(t, "GET", reqData.Meta.Method)
	assert.Equal(t, "/test?var1=foo&var2=bar", reqData.Meta.URL)
	assert.Equal(t, "foo", reqData.Input.Params["var1"])
	assert.Equal(t, "bar", reqData.Input.Params["var2"])
	assert.Equal(t, "Test Header", reqData.Meta.Headers["X-Test"])
}

func TestRequestDataPostWWWUrlEncoded(t *testing.T) {
	// Create a map of the form values
	data := url.Values{
		"name":  {Name},
		"email": {Email},
	}

	// Encode the form values
	encodedData := data.Encode()

	// Create the request
	req, err := http.NewRequest(
		"POST",
		"http://localhost/form",
		bytes.NewBufferString(encodedData),
	)
	require.NoError(t, err)

	// Add content-type header
	req.Header.Add(HeaderContentType, "application/x-www-form-urlencoded")

	reqData, err := requestdata.Collect(req, rules.Args{})
	require.NoError(t, err)

	assert.Equal(t, "POST", reqData.Meta.Method)
	assert.Equal(t, "/form", reqData.Meta.URL)
	assert.Equal(t, Name, reqData.Input.Form["name"])
	assert.Equal(t, Email, reqData.Input.Form["email"])
}

func TestRequestDataPostMultipartFormData(t *testing.T) {
	// Create buffer for body
	body := &bytes.Buffer{}

	// Create multipart writer
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("name", Name)
	_ = writer.WriteField("email", Email)

	// Add file
	file, err := os.Open("../../testdata/files/text.txt")
	require.NoError(t, err)
	part, err := writer.CreateFormFile("file", file.Name())
	require.NoError(t, err)
	_, err = io.Copy(part, file)
	require.NoError(t, err)

	// Close writer
	err = writer.Close()
	require.NoError(t, err)

	// Create request
	req, _ := http.NewRequest("POST", "http://localhost/upload", body)

	// Add headers
	req.Header.Add(HeaderContentType, writer.FormDataContentType())

	t.Run("with file uploads", func(t *testing.T) {
		reqData, err := requestdata.Collect(req, rules.Args{FileUploads: true})
		require.NoError(t, err)
		fStat, err := file.Stat()
		require.NoError(t, err)

		assert.Equal(t, "POST", reqData.Meta.Method)
		assert.Equal(t, "/upload", reqData.Meta.URL)
		assert.Equal(t, Name, reqData.Input.Form["name"])
		assert.Equal(t, Email, reqData.Input.Form["email"])
		assert.Equal(t, "file", reqData.Input.Uploads[0].FieldName)
		assert.Equal(t, "text/UTF-8", reqData.Input.Uploads[0].Type)
		assert.Equal(t, fStat.Size(), reqData.Input.Uploads[0].Size, "File size")
		assert.Contains(t, reqData.Input.Uploads[0].Stored, requestdata.UploadPrefix)
	})

	t.Run("without file uploads", func(t *testing.T) {
		reqData, err := requestdata.Collect(req, rules.Args{FileUploads: false})
		require.NoError(t, err)

		assert.Equal(t, "POST", reqData.Meta.Method)
		assert.Equal(t, "/upload", reqData.Meta.URL)
		assert.Equal(t, Name, reqData.Input.Form["name"])
		assert.Equal(t, Email, reqData.Input.Form["email"])
		assert.Equal(t, 0, len(reqData.Input.Uploads))
	})
}

func TestRequestDataPostJSON(t *testing.T) {
	// Create request body with JSON data
	body := bytes.NewBuffer([]byte(`{"Name":"John","Address":{"City":"London"}}`))

	// Create POST request
	req, err := http.NewRequest(
		"POST",
		"http://localhost/data",
		body,
	)
	require.NoError(t, err)

	// Set content type to JSON
	req.Header.Set(HeaderContentType, "application/json")

	reqData, err := requestdata.Collect(req, rules.Args{})
	require.NoError(t, err)

	assert.Equal(t, "POST", reqData.Meta.Method)
	assert.Equal(t, "/data", reqData.Meta.URL)
	assert.Equal(t, "John", extract.SFromI("Name", reqData.Input.JSON))
	assert.Equal(t, "London", extract.SFromI("Address.City", reqData.Input.JSON))
}

func TestURLPlaceholders(t *testing.T) {
	// Successful extraction with single placeholder
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	vars := map[string]string{
		"id": "foo",
	}

	req = mux.SetURLVars(req, vars)
	reqData, err := requestdata.Collect(req, rules.Args{})
	require.NoError(t, err)

	assert.Equal(t, "foo", reqData.Input.URLPlaceholders["id"])
}
