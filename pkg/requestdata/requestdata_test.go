package requestdata_test

import (
	"bytes"
	"encoding/json"
	"http-everything/httpe/pkg/requestdata"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const UserAgent = "Test User Agent"

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
	reqData, err := requestdata.Collect(req)
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
		"name":  {"John Doe"},
		"email": {"john@example.com"},
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
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	reqData, err := requestdata.Collect(req)
	require.NoError(t, err)

	assert.Equal(t, "POST", reqData.Meta.Method)
	assert.Equal(t, "/form", reqData.Meta.URL)
	assert.Equal(t, "John Doe", reqData.Input.Form["name"])
	assert.Equal(t, "john@example.com", reqData.Input.Form["email"])
}

func TestRequestDataPostMultipartFormData(t *testing.T) {
	// Create buffer for body
	body := &bytes.Buffer{}

	// Create multipart writer
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("name", "John Doe")
	_ = writer.WriteField("email", "john@example.com")

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
	req.Header.Add("Content-Type", writer.FormDataContentType())

	reqData, err := requestdata.Collect(req)
	require.NoError(t, err)

	assert.Equal(t, "POST", reqData.Meta.Method)
	assert.Equal(t, "/upload", reqData.Meta.URL)
	assert.Equal(t, "John Doe", reqData.Input.Form["name"])
	assert.Equal(t, "john@example.com", reqData.Input.Form["email"])
	assert.Equal(t, "file", reqData.Input.Uploads[0].FieldName)
	assert.Equal(t, "text/UTF-8", reqData.Input.Uploads[0].Type)
	assert.Equal(t, int64(21), reqData.Input.Uploads[0].Size)
	assert.Contains(t, reqData.Input.Uploads[0].Stored, requestdata.UploadPrefix)
}

func TestRequestDataPostJSON(t *testing.T) {
	// Create JSON data
	data := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
	}
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	// Create request body from JSON
	body := bytes.NewBuffer(jsonData)

	// Create POST request
	req, err := http.NewRequest(
		"POST",
		"http://localhost/data",
		body,
	)
	require.NoError(t, err)

	// Set content type to JSON
	req.Header.Set("Content-Type", "application/json")

	reqData, err := requestdata.Collect(req)
	require.NoError(t, err)

	assert.Equal(t, "POST", reqData.Meta.Method)
	assert.Equal(t, "/data", reqData.Meta.URL)
	t.Log(getValue(t, reqData.Input.JSON, "name"))
}

func TestURLPlaceholders(t *testing.T) {
	// Successful extraction with single placeholder
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	vars := map[string]string{
		"id": "foo",
	}

	req = mux.SetURLVars(req, vars)
	reqData, err := requestdata.Collect(req)
	require.NoError(t, err)

	assert.Equal(t, "foo", reqData.Input.URLPlaceholders["id"])
}

func getValue(t *testing.T, data interface{}, key string) string {
	t.Helper()

	// Handle nil data
	if data == nil {
		return ""
	}

	// Get the reflected Value
	val := reflect.ValueOf(data)

	// Handle invalid types
	if val.Kind() != reflect.Map {
		return ""
	}

	// Lookup key
	value := val.MapIndex(reflect.ValueOf(key))

	// Return value as interface
	return value.String()
}
