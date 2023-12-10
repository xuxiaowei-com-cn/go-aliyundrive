package aliyundrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://openapi.alipan.com/"
	userAgent      = "aliyundrive-go"
)

type Client struct {
	client *retryablehttp.Client

	baseURL *url.URL

	token string

	UserAgent string

	QrCode      *QrCodeService
	AccessToken *AccessTokenService
	UsersInfo   *UsersInfoService
}

type ListOptions struct {
	Page  int `url:"page,omitempty" json:"page,omitempty"`   // 当前的页码
	Limit int `url:"limit,omitempty" json:"limit,omitempty"` // 每页的数量
}

func NewClient(token string) (*Client, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	client.token = token
	return client, nil
}

// retryablehttp.NewClient()
func newClient() (*Client, error) {
	c := &Client{UserAgent: userAgent}

	c.client = &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultPooledClient(),
		Logger:       log.New(os.Stderr, "", log.LstdFlags),
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
		RetryMax:     4,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}

	err := c.SetBaseURL(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	c.QrCode = &QrCodeService{client: c}
	c.AccessToken = &AccessTokenService{client: c}
	c.UsersInfo = &UsersInfoService{client: c}

	return c, nil
}

func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

func (c *Client) SetBaseURL(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseURL = baseURL

	return nil
}

func (c *Client) NewRequest(method string, path string, requestQuery interface{}, requestBody interface{}) (*retryablehttp.Request, error) {

	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	var body interface{}
	switch {
	case method == http.MethodPatch || method == http.MethodPost || method == http.MethodPut:
		reqHeaders.Set("Content-Type", "application/json")

		if requestBody != nil {
			body, err = json.Marshal(requestBody)
			if err != nil {
				return nil, err
			}
		}
	}

	q, err := query.Values(requestQuery)
	if err != nil {
		return nil, err
	}

	u.RawQuery = q.Encode()

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}

	return req, nil
}

func (c *Client) UploadRequest(method, path string, content io.Reader, filename string) (*retryablehttp.Request, error) {
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(fw, content); err != nil {
		return nil, err
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	reqHeaders.Set("Content-Type", w.FormDataContentType())

	req, err := retryablehttp.NewRequest(method, u.String(), b)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}

	return req, nil
}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*Response, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

func CheckResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = fmt.Sprintf("failed to parse unknown error format: %s", data)
		} else {
			errorResponse.Message = parseError(raw)
		}
	}

	return errorResponse
}

func parseError(raw interface{}) string {
	switch raw := raw.(type) {
	case string:
		return raw

	case []interface{}:
		var errs []string
		for _, v := range raw {
			errs = append(errs, parseError(v))
		}
		return fmt.Sprintf("[%s]", strings.Join(errs, ", "))

	case map[string]interface{}:
		var errs []string
		for k, v := range raw {
			errs = append(errs, fmt.Sprintf("{%s: %s}", k, parseError(v)))
		}
		sort.Strings(errs)
		return strings.Join(errs, ", ")

	default:
		return fmt.Sprintf("failed to parse unexpected error type: %T", raw)
	}
}

func GetenvInt(key string, defaultValue int) int {
	var value = os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return i
}
