package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
    "path"
    "strconv"

	"github.com/pkg/errors"
)

type Node struct {
	endpointURL *url.URL
	httpClient  *http.Client
}

func NewNode(addr string) (*Node, error) {
	parsedURL, err := url.ParseRequestURI(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", addr)
	}

	n := &Node{
		endpointURL: parsedURL,
		httpClient:  &http.Client{},
	}
	return n, nil
}

func (node *Node) newRequest(ctx context.Context, method string, subPath string, body io.Reader) (*http.Request, error) {
	endpointURL := *node.endpointURL
	endpointURL.Path = path.Join(node.endpointURL.Path, subPath)

	req, err := http.NewRequest(method, endpointURL.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

type Status struct {
	Enable bool `json:"enable"`
}

func (node *Node) GetStatus(ctx context.Context) (*Status, error) {
	subPath := fmt.Sprintf("/api/status")
	req, err := node.newRequest(ctx, "GET", subPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := node.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var stat Status
	if err := decodeBody(resp, &stat); err != nil {
		return nil, err
	}

	return &stat, nil
}

func (node *Node) PostConfig(ctx context.Context, enable bool) error {
	subPath := fmt.Sprintf("/api/config")
	val := url.Values{}
	val.Set("enable", strconv.FormatBool(enable))
	req, err := node.newRequest(ctx, "POST", subPath, val)
	if err != nil {
		return err
	}

	resp, err := node.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return nil
}
