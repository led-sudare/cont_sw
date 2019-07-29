package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/pkg/errors"
)

type Node struct {
	endpointURL *url.URL
	httpClient  *http.Client
}

func NewNode(addr string) (*Node, error) {
	_addr := addProtocols(addr)
	parsedURL, err := url.ParseRequestURI(_addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", addr)
	}

	n := &Node{
		endpointURL: parsedURL,
		httpClient:  &http.Client{},
	}
	return n, nil
}

func addProtocols(addr string) string {
	if strings.HasPrefix(addr, "http://") ||
		strings.HasPrefix(addr, "https://") {
		return addr
	}
	return "http://" + addr
}

func (node *Node) newRequest(ctx context.Context, method string, subPath string, body io.Reader) (*http.Request, error) {
	endpointURL := *node.endpointURL
	endpointURL.Path = path.Join(node.endpointURL.Path, subPath)
	req, err := http.NewRequest(method, endpointURL.String(), body)
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}
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

	log.Debug(stat)
	return &stat, nil
}

func (node *Node) PostConfig(ctx context.Context, enable bool) error {
	subPath := fmt.Sprintf("/api/config")
	val, _ := json.Marshal(Status{Enable: enable})
	req, err := node.newRequest(ctx, "POST", subPath, bytes.NewBuffer(val))
	if err != nil {
		return err
	}

	dump1, _ := httputil.DumpRequestOut(req, true)

	resp, err := node.httpClient.Do(req)
	if err != nil {
		return err
	}

	dump2, _ := httputil.DumpResponse(resp, true)
	log.Debug(string(dump1))
	log.Debug(string(dump2))
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		log.Debug(resp.Body)
	}

	return nil
}
