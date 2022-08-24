package tinybird

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Request struct {
	Elapsed  string
	Method   string
	Pipe     Pipe
	Response Response
}

// Custom HTTP client for this module.
var Client HTTPClient

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Initialize module.
func init() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	Client = &http.Client{
		Timeout:   time.Duration(30) * time.Second,
		Transport: transport,
	}
}

func (r *Request) Execute() (err error) {
	r.Elapsed, err = Duration(func() error {
		req, err := r.newRequest()
		if err != nil {
			return err
		}

		res, err := Client.Do(req)
		if err != nil {
			return err
		}

		return r.readBody(res)
	})

	return err
}

func (r *Request) newRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.Pipe.GetURL(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.Pipe.Workspace.Token))
	req.URL.RawQuery = r.Pipe.Parameters.Encode()

	return req, nil
}

func (r *Request) readBody(resp *http.Response) (err error) {
	defer resp.Body.Close()

	r.Response.Status = resp.StatusCode
	r.Response.Body, err = io.ReadAll(resp.Body)
	r.Response.Decode()

	return err
}
