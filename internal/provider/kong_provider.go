package provider

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/kong/go-kong/kong"
	"github.com/pkg/errors"
)

// Config holds config details to use to talk to the Kong admin API.
type Config struct {
	Address            string
	Username           string
	Password           string
	InsecureSkipVerify bool
	APIKey             string
	AdminToken         string
	Workspace          string
}

// HeaderRoundTripper injects Headers into requests
// made via RT.
type HeaderRoundTripper struct {
	headers []string
	rt      http.RoundTripper
}

// IPPort represents a source and destination IP and port mapping used in routes.
type IPPort struct {
	IP   *string
	Port *int
}

// RoundTrip satisfies the RoundTripper interface.
func (t *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response,
	error) {
	newRequest := new(http.Request)
	*newRequest = *req
	newRequest.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		newRequest.Header[k] = append([]string(nil), s...)
	}
	for _, s := range t.headers {
		split := strings.SplitN(s, ":", 2)
		if len(split) >= 2 {
			newRequest.Header[split[0]] = append([]string(nil), split[1])
		}
	}
	return t.rt.RoundTrip(newRequest)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetKongClient(opt Config) (*kong.Client, error) {

	var tlsConfig tls.Config
	if opt.InsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	c := &http.Client{}
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("DefaultTransport is wrong type.")
	}
	defaultTransport.TLSClientConfig = &tlsConfig
	c.Transport = defaultTransport

	var headers []string
	if opt.APIKey != "" {
		headers = append(headers, fmt.Sprintf("apikey:%v", opt.APIKey))
	}
	if opt.AdminToken != "" {
		headers = append(headers, fmt.Sprintf("kong-admin-token:%v", opt.AdminToken))
	}
	if opt.Username != "" || opt.Password != "" {
		headers = append(headers, fmt.Sprintf("Authorization: Basic %v", basicAuth(opt.Username, opt.Password)))
	}
	if len(headers) > 0 {
		c.Transport = &HeaderRoundTripper{
			headers: headers,
			rt:      defaultTransport,
		}
	}

	url, err := url.Parse(opt.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse kong address")
	}
	if opt.Workspace != "" {
		url.Path = path.Join(url.Path, opt.Workspace)
	}

	kongClient, err := kong.NewClient(kong.String(url.String()), c)
	if err != nil {
		return nil, errors.Wrap(err, "creating client for Kong's Admin API")
	}

	return kongClient, nil
}
