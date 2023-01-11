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
)

type AddMemberCmd struct {
	Debug bool `help:"debugging output."`

	GheUrl string `arg name:"url" help:"url to Github Enterprise"`
	Token  string `arg name:"token" help:"token to authenticate to Github"`

	Member       string   `arg name:"user" help:"new member's info"`
	Organisation string   `arg name:"org" help:"owner or organisation for the new member"`
	Teams        []string `arg name:"teams" help:"the teams to add the new member to"`

	TlsConfig struct {
		HandshakeTimeout   int    `help:"Timeout in seconds" default:"10" env:"GHTOKEN_TIMEOUT"`
		InsecureSkipVerify bool   `help:"Allow insecure connections" default:"false" env:"GHTOKEN_INSECURE_SKIP_VERIFY"`
		Proxy              string `help:"Proxy URL" default:"" env:"GHTOKEN_PROXY"`
	} `embed:"" prefix:"tlsconfig."`
}

func (r *AddMemberCmd) Run() error {

	// Auth with GHE
	appAuth := AppAuthRequest{
		GheUrl:    r.GheUrl,
		Token:     r.Token,
		TlsConfig: r.TlsConfig,
	}

	client, appAuthErr := appAuth.Do()
	if appAuthErr != nil {
		return appAuthErr
	}

	// Request Merge
	collabreq := AddMemberRequest{
		Member:       r.Member,
		Organisation: r.Organisation,
		Teams:        r.Teams,
		Debug:        r.Debug,
	}
	err := collabreq.Do(context.TODO(), client)
	if err != nil {
		return err
	}

	return nil
}
