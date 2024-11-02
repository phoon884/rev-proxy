package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/phoon884/rev-proxy/internal/httpHdl/models"
)

func getMsgBeforeBody(buffer *bufio.Reader) ([]byte, error) {
	contentBuilder := bytes.Buffer{}
	for {
		line, err := buffer.ReadBytes('\n')
		if err == io.EOF || bytes.Equal(line, []byte{'\r', '\n'}) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			continue
		}
		contentBuilder.Write(line)
	}
	return contentBuilder.Bytes(), nil
}

func ParseReq(buffer *bufio.Reader) (*models.HTTPReq, error) {
	parsedMsg := models.HTTPReq{}
	content, err := getMsgBeforeBody(buffer)
	if len(content) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	header := make(map[string]string)
	headerLineSep := bytes.Split(content, []byte("\r\n"))
	for idx, line := range headerLineSep {
		if idx == 0 {
			elements := bytes.Split(line, []byte(" "))
			if len(elements) != 3 || string(elements[2]) != "HTTP/1.1" {
				return nil, errors.New("Malformed HTTP")
			}
			parsedMsg.Method = string(elements[0])
			pathAndParam := string(elements[1])
			pathAndParamSplice := strings.Split(pathAndParam, "?")
			switch len(pathAndParamSplice) {
			case 1:
				parsedMsg.Path = pathAndParamSplice[0]
				parsedMsg.Param = ""
			case 2:
				parsedMsg.Path = pathAndParamSplice[0]
				parsedMsg.Param = pathAndParamSplice[1]
			default:
				return nil, errors.New("Path error")
			}
		} else if bytes.Compare(line, []byte("\n")) == 0 {
		} else {
			headerKeyValue := bytes.SplitN(line, []byte(": "), 2)
			if len(headerKeyValue) == 2 {
				header[string(headerKeyValue[0])] = strings.TrimSpace(string(headerKeyValue[1]))
			} else {
				if string(headerKeyValue[0]) != "" {
					header[string(headerKeyValue[0])] = ""
				}
			}
		}
	}
	parsedMsg.Header = header
	return &parsedMsg, nil
}

func ParseRes(buffer *bufio.Reader) (*models.HTTPRes, error) {
	parsedMsg := models.HTTPRes{}
	content, err := getMsgBeforeBody(buffer)
	if len(content) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	header := make(map[string]string)
	headerLineSep := bytes.Split(content, []byte("\r\n"))
	for idx, line := range headerLineSep {
		if idx == 0 {
			elements := bytes.SplitN(line, []byte(" "), 3)
			if string(elements[0]) != "HTTP/1.1" {
				return nil, errors.New("Malformed HTTP")
			}
			parsedMsg.ResponseCode, err = strconv.Atoi(string(elements[1]))
			if err != nil {
				return nil, err
			}
		} else if bytes.Compare(line, []byte("\n")) != 0 {
			headerKeyValue := bytes.SplitN(line, []byte(": "), 2)
			if len(headerKeyValue) == 2 {
				header[string(headerKeyValue[0])] = strings.TrimSpace(string(headerKeyValue[1]))
			} else {
				if string(headerKeyValue[0]) != "" {
					header[string(headerKeyValue[0])] = ""
				}
			}
		}
	}
	parsedMsg.Header = header

	fmt.Println("error non")
	return &parsedMsg, nil
}
