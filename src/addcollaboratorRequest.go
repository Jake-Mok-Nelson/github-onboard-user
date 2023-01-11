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
	"fmt"

	"github.com/google/go-github/v48/github"
)

// AddMemberRequest contains configuration required to onboard a customer to a Github org and teams
type AddMemberRequest struct {
	Member       string   // The new user's info
	Organisation string   // The organisation to target
	Teams        []string // The teams to add the user to
	Debug        bool
}

func (req AddMemberRequest) Do(ctx context.Context, client *github.Client) error {

	// Check that we can access the organisation
	_, _, err := client.Organizations.Get(ctx, req.Organisation)
	if err != nil {
		return fmt.Errorf("unable to access organisation %v", req.Organisation)
	}

	// For each team in req.Teams
	for _, team := range req.Teams {
		// Add the user to the team
		_, _, err := client.Teams.AddTeamMembershipBySlug(ctx, req.Organisation, team, req.Member, nil)

		if err != nil {
			return err
		}
		fmt.Printf("\nAttempting add user for team %v", team)
	}

	return nil
}
