package models

import (
	"bytes"
	"compress/gzip"
	. "github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/core/views"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

func TestConvertToResponseDetailsView_WithPlainTextResponseDetails(t *testing.T) {
	RegisterTestingT(t)

	statusCode := 200
	body := "hello_world"
	headers := map[string][]string{"test_header": []string{"true"}}

	originalResp := ResponseDetails{Status: statusCode, Body: body, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(false))
	Expect(respView.Body).To(Equal(body))
}

func TestConvertToResponseDetailsView_WithGzipContentEncodedHeader(t *testing.T) {
	RegisterTestingT(t)

	originalBody := "hello_world"

	statusCode := 200
	body := GzipString(originalBody)
	headers := map[string][]string{"Content-Encoding": []string{"gzip"}}

	originalResp := ResponseDetails{Status: statusCode, Body: body, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(true))
	Expect(respView.Body).NotTo(Equal(body))
	Expect(respView.Body).NotTo(Equal(originalBody))

	base64EncodedBody := "H4sIAAAJbogA/w=="

	Expect(respView.Body).To(Equal(base64EncodedBody))
}

func TestConvertToResponseDetailsView_WithDeflateContentEncodedHeader(t *testing.T) {
	RegisterTestingT(t)

	originalBody := "this_should_be_encoded_but_its_not_important"

	statusCode := 200
	headers := map[string][]string{"Content-Encoding": []string{"deflate"}}

	originalResp := ResponseDetails{Status: statusCode, Body: originalBody, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(true))
	Expect(respView.Body).NotTo(Equal(originalBody))

	base64EncodedBody := "dGhpc19zaG91bGRfYmVfZW5jb2RlZF9idXRfaXRzX25vdF9pbXBvcnRhbnQ="

	Expect(respView.Body).To(Equal(base64EncodedBody))
}

func TestConvertToResponseDetailsView_WithImageBody(t *testing.T) {
	RegisterTestingT(t)

	imageUri := "/testdata/1x1.png"

	file, _ := os.Open("../../functional-tests/core" + imageUri)
	defer file.Close()

	originalImageBytes, _ := ioutil.ReadAll(file)

	originalResp := ResponseDetails{
		Status: 200,
		Body:   string(originalImageBytes),
	}

	respView := originalResp.ConvertToResponseDetailsView()

	base64EncodedBody := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVR4nGP6DwABBQECz6AuzQAAAABJRU5ErkJggg=="
	Expect(respView).To(Equal(views.ResponseDetailsView{
		Status:      200,
		Body:        base64EncodedBody,
		EncodedBody: true,
	}))
}
func TestRequestResponsePair_ConvertToRequestResponsePairView_WithPlainTextResponse(t *testing.T) {
	RegisterTestingT(t)

	respBody := "hello_world"

	requestResponsePair := RequestResponsePair{
		Response: ResponseDetails{
			Status:  200,
			Body:    respBody,
			Headers: map[string][]string{"test_header": []string{"true"}}},
		Request: RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "",
			Body:        "",
			Headers:     map[string][]string{"test_header": []string{"true"}}},
	}

	pairView := requestResponsePair.ConvertToRequestResponsePairView()

	Expect(*pairView).To(Equal(views.RequestResponsePairView{
		Response: views.ResponseDetailsView{
			Status:      200,
			Body:        respBody,
			Headers:     map[string][]string{"test_header": []string{"true"}},
			EncodedBody: false},
		Request: views.RequestDetailsView{
			RequestType: StringToPointer("recording"),
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"test_header": []string{"true"}}},
	}))
}

func TestRequestResponsePair_ConvertToRequestResponsePairView_WithGzippedResponse(t *testing.T) {
	RegisterTestingT(t)

	requestResponsePair := RequestResponsePair{
		Response: ResponseDetails{
			Status:  200,
			Body:    GzipString("hello_world"),
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "",
			Body:        "",
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}},
		},
	}

	pairView := requestResponsePair.ConvertToRequestResponsePairView()

	Expect(*pairView).To(Equal(views.RequestResponsePairView{
		Response: views.ResponseDetailsView{
			Status:      200,
			Body:        "H4sIAAAJbogA/w==",
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}},
			EncodedBody: true},
		Request: views.RequestDetailsView{
			RequestType: StringToPointer("recording"),
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}},
		},
	}))
}

func TestRequestDetails_ConvertToRequestDetailsView(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := RequestDetails{
		Path:        "/",
		Method:      "GET",
		Destination: "/",
		Scheme:      "scheme",
		Query:       "", Body: "",
		Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}}

	requestDetailsView := requestDetails.ConvertToRequestDetailsView()

	Expect(requestDetailsView.Path).To(Equal(StringToPointer(requestDetails.Path)))
	Expect(requestDetailsView.Method).To(Equal(StringToPointer(requestDetails.Method)))
	Expect(requestDetailsView.Destination).To(Equal(StringToPointer(requestDetails.Destination)))
	Expect(requestDetailsView.Scheme).To(Equal(StringToPointer(requestDetails.Scheme)))
	Expect(requestDetailsView.Query).To(Equal(StringToPointer(requestDetails.Query)))
	Expect(requestDetailsView.Headers).To(Equal(requestDetails.Headers))
}

// Helper function for gzipping strings
func GzipString(s string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}

func TestRequestResponsePairView_ConvertToRequestResponsePairWithoutEncoding(t *testing.T) {
	RegisterTestingT(t)

	view := views.RequestResponsePairView{
		Request: views.RequestDetailsView{
			Path:        StringToPointer("A"),
			Method:      StringToPointer("A"),
			Destination: StringToPointer("A"),
			Scheme:      StringToPointer("A"),
			Query:       StringToPointer("A"),
			Body:        StringToPointer("A"),
			Headers: map[string][]string{
				"A": []string{"B"},
				"C": []string{"D"},
			},
		},
		Response: views.ResponseDetailsView{
			Status:      1,
			Body:        "1",
			EncodedBody: false,
			Headers: map[string][]string{
				"1": []string{"2"},
				"3": []string{"4"},
			},
		},
	}

	requestResponsePair := NewRequestResponsePairFromRequestResponsePairView(view)

	Expect(requestResponsePair).To(Equal(RequestResponsePair{
		Request: RequestDetails{
			Path:        "A",
			Method:      "A",
			Destination: "A",
			Scheme:      "A",
			Query:       "A",
			Body:        "A",
			Headers: map[string][]string{
				"A": []string{"B"},
				"C": []string{"D"},
			},
		},
		Response: ResponseDetails{
			Status: 1,
			Body:   "1",
			Headers: map[string][]string{
				"1": []string{"2"},
				"3": []string{"4"},
			},
		},
	}))
}

func TestRequestResponsePairView_ConvertToRequestResponsePairWithEncoding(t *testing.T) {
	RegisterTestingT(t)

	view := views.RequestResponsePairView{
		Response: views.ResponseDetailsView{
			Body:        "ZW5jb2RlZA==",
			EncodedBody: true,
		},
	}

	pair := NewRequestResponsePairFromRequestResponsePairView(view)

	Expect(pair.Response.Body).To(Equal("encoded"))
}

func TestRequestDetailsView_ConvertToRequestDetails(t *testing.T) {
	RegisterTestingT(t)

	requestDetailsView := views.RequestDetailsView{
		Path:        StringToPointer("/"),
		Method:      StringToPointer("GET"),
		Destination: StringToPointer("/"),
		Scheme:      StringToPointer("scheme"),
		Query:       StringToPointer(""),
		Body:        StringToPointer(""),
		Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}}}

	requestDetails := NewRequestDetailsFromRequestDetailsView(requestDetailsView)

	Expect(requestDetails.Path).To(Equal(*requestDetailsView.Path))
	Expect(requestDetails.Method).To(Equal(*requestDetailsView.Method))
	Expect(requestDetails.Destination).To(Equal(*requestDetailsView.Destination))
	Expect(requestDetails.Scheme).To(Equal(*requestDetailsView.Scheme))
	Expect(requestDetails.Query).To(Equal(*requestDetailsView.Query))
	Expect(requestDetails.Headers).To(Equal(requestDetailsView.Headers))
}
