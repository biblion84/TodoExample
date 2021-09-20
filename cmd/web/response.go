package main

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	http.Response
	stringBody string
}

func (res *Response) String() string {
	if res.stringBody != "" {
		return res.stringBody
	}
	if res == nil || res.Body == nil {
		return ""
	} else {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	res.stringBody = string(body)
	return res.stringBody
}
func (res *Response) Ok() bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

func GetResponse(response *http.Response, err error) (Response, error) {
	if err != nil {
		return Response{}, err
	}
	return Response{*response, ""}, nil
}
