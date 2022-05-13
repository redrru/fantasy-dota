package http

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	client *retryablehttp.Client
}

func NewClient(logger *log.Logger) *Client {
	reClient := retryablehttp.NewClient()
	reClient.RetryMax = 20
	reClient.RetryWaitMin = 10 * time.Millisecond
	reClient.RetryWaitMax = 1 * time.Second
	reClient.Logger = logger

	return &Client{
		client: reClient,
	}
}

func (c *Client) Get(url string) ([]byte, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (c *Client) Post(url string, body []byte) error {
	req, err := retryablehttp.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, _ = io.Copy(ioutil.Discard, res.Body)

	return nil
}
