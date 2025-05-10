package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	url    string
	method string
	header Header
	param  Param
	body   Body
}

func NewRequest(c Client, url, method string, needCredential bool) *Request {
	req := Request{
		url:    fmt.Sprintf("%s/%s", c.base, url),
		method: method,
		header: Header{
			"content-type": "application/json",
		},
		param: Param{},
		body:  Body{},
	}

	if needCredential {
		if c.token != "" {
			req.header["authorization"] = fmt.Sprintf("bearer %s", c.token)
		}
		req.header["appKey"] = c.appKey
		req.header["appSecret"] = c.appSecret
	}

	return &req
}

func (r *Request) SetHeader(key, content string) {
	r.header[key] = content
}

func (r *Request) SetParam(key, content string) {
	r.param[key] = content
}

func (r *Request) SetBody(key string, content interface{}) {
	r.body[key] = content
}

func (r Request) Parse() (*http.Request, error) {
	u, err := url.Parse(r.url)
	if err != nil {
		return nil, fmt.Errorf("error while parsing %s:%v", r.url, err)
	}

	q := u.Query()
	for key, content := range r.param {
		q.Set(key, content)
	}
	u.RawQuery = q.Encode()

	b, err := r.body.Marshal()
	if err != nil {
		return nil, fmt.Errorf("error while marshallig body:%v", err)
	}

	req, _ := http.NewRequest(r.method, u.String(), bytes.NewBuffer(b))
	for key, content := range r.header {
		req.Header.Set(key, content)
	}

	return req, nil
}

func (r *Request) Send(result interface{}) error {
	req, err := r.Parse()
	if err != nil {
		return fmt.Errorf("error while parsing request:%v", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while doing request:%v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error while reading response (%v):%v", res, err)
	}

	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("failed to unmarshal:%v", err)
	}

	return nil
}
