package artifactory

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"strings"
)

type HttpsClient struct {
	Client *http.Client
	req    *http.Request
	pool   *x509.CertPool
}

func (c *HttpsClient) BasicAuth(user, pass string) Client {
	c.req = &http.Request{}
	c.req.Header = make(map[string][]string)
	c.req.SetBasicAuth(user, pass)
	c.req.Header.Set("Content-Type", "application/json")
	return c
}

func (c *HttpsClient) Certificate(certpath string) Client {
	c.pool = x509.NewCertPool()
	caCrt, _ := ioutil.ReadFile(certpath)
	c.pool.AppendCertsFromPEM(caCrt)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: c.pool,
			//TODO: 关闭Insecure
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}

	c.Client = &http.Client{Transport: tr}
	return c
}

func (c *HttpsClient) Get(path string) (int, string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return 0, "", err
	}

	c.req.URL = u
	c.req.Method = http.MethodGet
	code, resp, err := c.do(c.req)
	data, err := getResponseData(resp)
	if err != nil {
		return 0, "", err
	}

	return code, data, nil
}

type HttpsClientReader struct {
	*bytes.Reader
}

// Read(p []byte) (n int, err error)
func (r *HttpsClientReader)Close() error {
	return nil
}


func (c *HttpsClient) Post(path string, body []byte) (int, string, error) {
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	if user, pass, ok := c.req.BasicAuth(); ok {
		req.SetBasicAuth(user, pass)
	}
	if ok := c.req.Header.Get("Content-Type") != ""; ok {
		req.Header.Set("Content-Type", c.req.Header.Get("Content-Type"))
	}
	c.req = req

	code, resp, err := c.do(c.req)
	if err != nil {
		return 0, "", err
	}

	data, err := getResponseData(resp)
	if err != nil {
		return 0, "", err
	}

	return code, data, nil
}

func (c *HttpsClient) Delete(path string) (int, string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return 0, "", err
	}

	c.req.URL = u
	c.req.Method = http.MethodDelete

	resp, err := c.Client.Do(c.req)
	if err != nil {
		return 0, "", err
	}

	data, err := getResponseData(resp)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, data, nil
}

func (c *HttpsClient) Put(path string, body []byte) (int, string, error) {
	req, err := http.NewRequest(http.MethodPut, path, strings.NewReader(string(body)))
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	if user, pass, ok := c.req.BasicAuth(); ok {
		req.SetBasicAuth(user, pass)
	}
	if ok := c.req.Header.Get("Content-Type") != ""; ok {
		req.Header.Set("Content-Type", c.req.Header.Get("Content-Type"))
	}
	c.req = req

	code, resp, err := c.do(c.req)
	if err != nil {
		return 0, "", err
	}

	data, err := getResponseData(resp)
	if err != nil {
		return 0, "", err
	}

	return code, data, nil
}

func (c *HttpsClient) do(req *http.Request) (int, *http.Response, error) {
	resp, err := c.Client.Do(c.req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return resp.StatusCode, resp, err
}

func getResponseData(req *http.Response) (string, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *HttpsClient)RawGet(path string)(int, *http.Response, error)  {
	u, err := url.Parse(path)
	if err != nil {
		return 0, nil, err
	}

	c.req.URL = u
	c.req.Method = http.MethodGet
	code, resp, err := c.do(c.req)
	return code, resp, nil
}
