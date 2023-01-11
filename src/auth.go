/*
Copyright 2021 Jake Nelson

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute,
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type AppAuthRequest struct {
	GheUrl string // Github Enterprise URL
	Token  string

	TlsConfig struct {
		HandshakeTimeout   int    `help:"Timeout in seconds" default:"10" env:"GHTOKEN_TIMEOUT"`
		InsecureSkipVerify bool   `help:"Allow insecure connections" default:"false" env:"GHTOKEN_INSECURE_SKIP_VERIFY"`
		Proxy              string `help:"Proxy URL" default:"" env:"GHTOKEN_PROXY"`
	} `embed:"" prefix:"tlsconfig."`
}

func (config AppAuthRequest) Do() (*github.Client, error) {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	httpClient := oauth2.NewClient(ctx, tokenSource)

	httpClient.Timeout = time.Duration(config.TlsConfig.HandshakeTimeout) * time.Second
	httpClient.Transport = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.TlsConfig.InsecureSkipVerify},
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: time.Duration(config.TlsConfig.HandshakeTimeout) * time.Second,
	}

	client, err := github.NewEnterpriseClient(config.GheUrl, config.GheUrl, httpClient)
	if err != nil {
		return nil, err
	}

	return client, nil
}
