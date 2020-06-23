package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tcz001/databricks-sdk-go/models"
	"golang.org/x/time/rate"
)

type Options struct {
	Domain *string
	Token  *string

	XDatabricksAzureWorkspaceResourceId *string
	XDatabricksAzureSPManagementToken   *string

	MaxRetries         int
	RetryDelay         time.Duration
	RateLimitPerSecond int
}

type Client struct {
	http        *http.Client
	baseUrl     *url.URL
	header      http.Header
	maxRetries  int
	retryDelay  time.Duration
	rateLimiter *rate.Limiter
}

func NewClient(opts Options) (*Client, error) {
	loadEnvConfig(&opts)

	if opts.Domain == nil {
		return nil, fmt.Errorf("missing Domain")
	}

	if opts.XDatabricksAzureWorkspaceResourceId != nil && (opts.Token == nil || opts.XDatabricksAzureSPManagementToken == nil) {
		return nil, fmt.Errorf("missing X-Databricks-Azure-SP-Management-Token, when specifiying X-Databricks-Azure-Workspace-Resource-Id is provided")
	}

	if opts.XDatabricksAzureSPManagementToken != nil && opts.XDatabricksAzureWorkspaceResourceId == nil {
		return nil, fmt.Errorf("missing X-Databricks-Azure-Workspace-Resource-Id, when X-Databricks-Azure-SP-Management-Token is provided")
	}

	baseUrl, err := url.Parse(fmt.Sprintf("https://%s/api/2.0/", *opts.Domain))
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	if opts.Token != nil {
		header.Add("Authorization", fmt.Sprintf("Bearer %s", *opts.Token))
	}
	if opts.XDatabricksAzureWorkspaceResourceId != nil && opts.XDatabricksAzureSPManagementToken != nil {
		header.Add("X-Databricks-Azure-SP-Management-Token", fmt.Sprintf(*opts.XDatabricksAzureSPManagementToken))
		header.Add("X-Databricks-Azure-Workspace-Resource-Id", fmt.Sprintf(*opts.XDatabricksAzureWorkspaceResourceId))
	}

	client := Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		baseUrl:     baseUrl,
		header:      header,
		maxRetries:  opts.MaxRetries,
		retryDelay:  opts.RetryDelay,
		rateLimiter: rate.NewLimiter(rate.Limit(opts.RateLimitPerSecond), 1),
	}

	return &client, nil
}

func loadEnvConfig(opts *Options) {
	if opts.Domain == nil {
		if v := os.Getenv("DATABRICKS_DOMAIN"); v != "" {
			opts.Domain = &v
		}
	}

	if opts.Token == nil {
		if v := os.Getenv("DATABRICKS_TOKEN"); v != "" {
			opts.Token = &v
		}
	}
}

func (c *Client) Query(method string, path string, data interface{}) ([]byte, error) {
	request, err := c.buildRequest(method, path, data)
	if err != nil {
		return nil, err
	}

	var responseBytes []byte

	for i := 0; ; i++ {
		err := c.rateLimiter.Wait(context.TODO())
		if err == nil {
			fmt.Printf("rate limit error: %v", err)
		}
		responseBytes, err = c.makeRequest(request)
		if err == nil {
			break
		}

		temporary := isTemporary(err)
		if !temporary || i >= c.maxRetries {
			break
		}

		time.Sleep(c.retryDelay)
	}

	return responseBytes, err
}

func (c *Client) buildRequest(method string, path string, data interface{}) (*http.Request, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	queryUrl := c.baseUrl.ResolveReference(u)

	var body []byte = nil
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest(method, queryUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header = c.header

	return request, nil
}

func (c *Client) makeRequest(request *http.Request) ([]byte, error) {
	log.Printf("[DEBUG] HTTP request: %v", request)

	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] HTTP response: %v", response)

	defer response.Body.Close()

	return c.parseResponse(request, *response)
}

func (c *Client) parseResponse(request *http.Request, response http.Response) ([]byte, error) {
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Response bytes: %s", responseBytes)

	if response.StatusCode == 204 {
		return nil, nil
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		log.Printf("[ERROR] Error Response StatusCode: %d", response.StatusCode)
		errorResponse := models.ErrorResponse{}

		if strings.Contains(response.Header.Get("Content-Type"), "json") {
			err = json.Unmarshal(responseBytes, &errorResponse)
			if err != nil {
				log.Printf("[ERROR] Error Response : %v", errorResponse)
				return nil, err
			}
		} else {
			errorResponse.Message = fmt.Sprintf(
				"databricks request: %v %#v error: %#v %s", request.URL, request, response, string(responseBytes))
			log.Printf("[ERROR] Error Response Message: %s", errorResponse.Message)
		}

		return nil, NewError(errorResponse, response.StatusCode)
	}

	return responseBytes, nil
}

func isTemporary(err error) bool {
	if nerr, ok := err.(net.Error); ok {
		return nerr.Temporary()
	}

	if derr, ok := err.(Error); ok {
		return derr.Temporary()
	}

	return false
}

type Error struct {
	errorResponse models.ErrorResponse
	statusCode    int
}

func NewError(response models.ErrorResponse, statusCode int) Error {
	return Error{
		errorResponse: response,
		statusCode:    statusCode,
	}
}

func (e Error) Error() string {
	return e.errorResponse.Message
}

func (e Error) Code() string {
	return e.errorResponse.ErrorCode
}

func (e Error) Temporary() bool {
	return e.statusCode >= 500
}
