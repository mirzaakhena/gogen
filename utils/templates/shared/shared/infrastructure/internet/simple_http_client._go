package internet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SimpleHTTPClient struct {
	url     string
	method  string
	body    io.Reader
	errs    []error
	headers map[string]string
	timeout time.Duration
}

// NewSimpleHTTPClient is the fastest way to use the http client for calling any API by default it already have header
//
//	"Content-Type": "application/json"
//
// and the timeout is 3 second
func NewSimpleHTTPClient(method, url string, requestData ...any) *SimpleHTTPClient {

	errs := make([]error, 0)

	var err error
	var body io.Reader

	if len(requestData) > 0 {
		body, err = constructBody(requestData[0])
		if err != nil {
			errs = append(errs, err)
		}
	}

	return &SimpleHTTPClient{
		url:    url,
		method: method,
		body:   body,
		errs:   errs,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		timeout: time.Second * 3,
	}
}

func (r *SimpleHTTPClient) Method(method string) *SimpleHTTPClient {
	if r.method != "" {
		return r
	}
	r.method = method
	return r
}

func (r *SimpleHTTPClient) URL(url string) *SimpleHTTPClient {
	if r.url != "" {
		return r
	}
	if url == "" {
		r.errs = append(r.errs, fmt.Errorf("url must not empty"))
	}
	r.url = url
	return r
}

func constructBody(requestData any) (io.Reader, error) {
	jsonInBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonInBytes), nil
}

func (r *SimpleHTTPClient) Body(requestData any) *SimpleHTTPClient {
	if r.body != nil {
		return r
	}

	body, err := constructBody(requestData)
	if err != nil {
		r.errs = append(r.errs, err)
		return r
	}

	r.body = body

	return r
}

func (r *SimpleHTTPClient) Header(key string, value string) *SimpleHTTPClient {
	r.headers[key] = value
	return r
}

func (r *SimpleHTTPClient) Timeout(duration time.Duration) *SimpleHTTPClient {
	r.timeout = duration
	return r
}

type ResponseData struct {
	Data    any
	Code    int
	Status  string
	Headers map[string][]string
}

func (r *SimpleHTTPClient) Call(responseData *ResponseData) error {

	if len(r.errs) > 0 {
		errMessage := ""
		for i, s := range r.errs {
			if i == 0 {
				errMessage += s.Error()
				continue
			}
			errMessage += fmt.Sprintf(", %s", s.Error())
		}
		return fmt.Errorf(errMessage)
	}

	request, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return err
	}

	for k, v := range r.headers {
		request.Header.Set(k, v)
	}

	var client = &http.Client{Timeout: r.timeout}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, responseData)
	if err != nil {
		return err
	}

	responseData = &ResponseData{
		Data:    responseData,
		Code:    response.StatusCode,
		Status:  response.Status,
		Headers: response.Header,
	}

	return nil
}
