package models

import (
	"fmt"
	"strings"
)

var stringPresentableContentTypes = map[string]struct{}{
	"text/plain":             {},
	"text/html":              {},
	"text/css":               {},
	"text/csv":               {},
	"text/javascript":        {},
	"application/json":       {},
	"application/xml":        {},
	"application/xhtml+xml":  {},
	"application/javascript": {},
	"application/x-yaml":     {},
	"application/rtf":        {},
	"text/markdown":          {},
	"text/event-stream":      {},
}

type HTTPReq struct {
	Method string
	Path   string
	Param  string
	Header map[string]string
	Body   []byte
}

// Returns a string representation of the request will not show body if Content-Type is not a text
func (h HTTPReq) StrRepr() string {
	strRepr := strings.Builder{}
	title := fmt.Sprintf("%s %s HTTP/1.1\r\n", h.Method, h.Path)
	strRepr.WriteString(title)
	for key, value := range h.Header {
		strRepr.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	strRepr.WriteString("\r\n")
	return strRepr.String()
}

func (h *HTTPReq) ChangeHeader(header string, value string) {
	if h.Header == nil {
		h.Header = make(map[string]string)
	}
	h.Header[header] = value
}

func (h HTTPReq) GetHeader(header string) string {
	result := ""
	if h.Header != nil {
		result = h.Header[header]
	}
	return result
}
