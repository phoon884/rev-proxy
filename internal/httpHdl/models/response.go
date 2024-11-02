package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

var httpStatusText = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

type HTTPRes struct {
	ResponseCode int
	Header       map[string]string
	Body         []byte
}

func (r HTTPRes) ToBytes() []byte {
	res := bytes.Buffer{}
	responseText, ok := httpStatusText[r.ResponseCode]
	if !ok {
		responseText = "N/A"
	}
	firstLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.ResponseCode, responseText)
	res.WriteString(firstLine)
	for key, value := range r.Header {
		res.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	res.WriteString("\r\n")
	if len(r.Body) != 0 {
		res.WriteString("\r\n")
		res.Write(r.Body)
	}
	return res.Bytes()
}

func NewErrorFound(msgCode int, msg string) HTTPRes {
	body := []byte(fmt.Sprintf(`{"Error":"%s"}`, msg))
	res := HTTPRes{}
	header := map[string]string{
		"Date":            time.Now().String(),
		"Server":          "Go-Reverse-Proxy",
		"Content-Type":    "application/json",
		"Conntent-Length": fmt.Sprint(binary.Size(body)),
		"Connection":      "close",
	}
	res.Header = header
	res.Body = body
	res.ResponseCode = msgCode
	return res
}

func (h HTTPRes) GetHeader(header string) string {
	result := ""
	if h.Header != nil {
		result = h.Header[header]
	}
	return result
}
