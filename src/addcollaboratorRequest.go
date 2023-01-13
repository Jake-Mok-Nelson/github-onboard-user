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
	"net/http"

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

	// Get the organisation
	if req.Debug {
		fmt.Printf("Getting organisation %v\n", req.Organisation)
	}
	_, resp, err := client.Organizations.Get(ctx, req.Organisation)
	var targetOrg string
	if resp != nil && req.Debug {
		fmt.Printf("Response from Get Oragnization: %v\n", resp.Status)
	}
	if err != nil {
		return fmt.Errorf("unable to read organisation %v, err: %v", req.Organisation, err)
	}

	// Check if the user is in the org
	if req.Debug {
		fmt.Printf("Checking if %v is a member of %v\n", req.Member, req.Organisation)
	}
	_, resp, err = client.Organizations.GetOrgMembership(ctx, req.Member, req.Organisation)
	if resp != nil && req.Debug {
		fmt.Printf("Response from GetOrgMembership: %v\n", resp.Status)
	}

	// If getOrgContinue is true, we'll continue to add the user to the org
	var getOrgContinue bool
	if err != nil {
		// Check for StatusNotModified or StatusNotFound
		if resp.StatusCode != http.StatusNotModified && req.Debug {
			fmt.Printf("GetOrgMembership returned StatusNotModified, this probably means no change is required\n")
		}

		if resp.StatusCode != http.StatusNotFound && req.Debug {
			fmt.Printf("GetOrgMembership StatusNotFound: %v\n", err)
			getOrgContinue = true
		}
	}

	// If the user is not in the org, add them
	if getOrgContinue {
		newMembership, resp, err := client.Organizations.EditOrgMembership(ctx, req.Member, targetOrg, &github.Membership{})
		if resp != nil && req.Debug {
			fmt.Printf("Response from EditOrgMembership: %v\n", resp.Status)
		}
		if err != nil {
			return fmt.Errorf("unable to add user %v to organisation %v, err: %v", req.Member, req.Organisation, err)
		}

		stateOfNewUser := newMembership.GetState()
		if stateOfNewUser != "active" {
			return fmt.Errorf("unable to add user %v to organisation %v, err: %v", req.Member, req.Organisation, err)
		}
	}

	// If we reach here, we should be a member of the organisation
	// Attempt to add the user to each of the teams provided
	for _, team := range req.Teams {
		// Check if they are in the team
		_, resp, err := client.Teams.GetTeamMembershipBySlug(ctx, req.Organisation, team, req.Member)
		if resp != nil && req.Debug {
			fmt.Printf("Response from GetTeamBySlug: %v\n", resp.Status)
		}
		if err != nil {
			if resp != nil {
				if resp.StatusCode == http.StatusNotFound && req.Debug {
					fmt.Printf("User %v not in team %v\n", req.Member, team)
				}
				if resp.StatusCode == http.StatusNotModified && req.Debug {
					fmt.Printf("User %v already in team %v\n", req.Member, team)
					continue
				}
			}
			// return an error if it's not a 404 (if it's a 404, they aren't a member and we should add them)
			if resp != nil && resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("unable to read team %v, err: %v", team, err)
			}
			return fmt.Errorf("unable to read team %v", team)
		}

		// Add the user to the team
		_, resp, err = client.Teams.AddTeamMembershipBySlug(ctx, req.Organisation, team, req.Member, nil)
		if resp != nil && req.Debug {
			fmt.Printf("Response from AddTeamMembershipBySlug: %v\n", resp.Status)
		}
		// Check for some known statuses and handle them
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("unable to add user %v to team %v, err: %v", req.Member, team, err)
			}
			if resp.StatusCode == http.StatusNotModified {
				fmt.Printf("User %v already in team %v\n", req.Member, team)
				continue
			}

			return err
		}

		fmt.Printf("User %v added to team %v\n", req.Member, team)
	}

	return nil
}
