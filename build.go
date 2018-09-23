// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package travis

import (
	"context"
	"fmt"
	"net/http"
)

// BuildService handles communication with the builds
// related methods of the Travis CI API.
type BuildService struct {
	client *Client
}

// Build represents a Travis CI build
//
// Travis CI API docs: https://developer.travis-ci.com/resource/build#standard-representation
type Build struct {
	// Value uniquely identifying the build
	Id uint `json:"id,omitempty"`
	// Incremental number for a repository's builds
	Number string `json:"number,omitempty"`
	// Current state of the build
	State string `json:"state,omitempty"`
	// Wall clock time in seconds
	Duration uint `json:"duration,omitempty"`
	// Event that triggered the build
	EventType string `json:"event_type,omitempty"`
	// State of the previous build (useful to see if state changed)
	PreviousState string `json:"previous_state,omitempty"`
	// Title of the build's pull request
	PullRequestTitle string `json:"pull_request_title,omitempty"`
	// Number of the build's pull request
	PullRequestNumber uint `json:"pull_request_number,omitempty"`
	// When the build started
	StartedAt string `json:"started_at,omitempty"`
	// When the build finished
	FinishedAt string `json:"finished_at,omitempty"`
	// The last time the build was updated
	UpdatedAt string `json:"updated_at,omitempty"`
	// Whether or not the build is private
	Private bool `json:"private,omitempty"`
	// GitHub repository the build is associated with
	Repository MinimalRepository `json:"repository,omitempty"`
	// The branch the build is associated with
	Branch MinimalBranch `json:"branch,omitempty"`
	// The build's tag
	Tag string `json:"tag,omitempty"`
	// The commit the build is associated with
	Commit MinimalCommit `json:"commit,omitempty"`
	// List of jobs that are part of the build's matrix
	Jobs []MinimalJob `json:"jobs,omitempty"`
	// The stages of the build
	Stages []MinimalStage `json:"stages,omitempty"`
	// The User or Organization that created the build
	CreatedBy MinimalOwner `json:"owner,omitempty"`
}

// MinimalBuild is a minimal representation of a Travis CI build
//
// Travis CI API docs: https://developer.travis-ci.com/resource/build#minimal-representation
type MinimalBuild struct {
	// Value uniquely identifying the build
	Id uint `json:"id,omitempty"`
	// Incremental number for a repository's builds
	Number string `json:"number,omitempty"`
	// Current state of the build
	State string `json:"state,omitempty"`
	// Wall clock time in seconds
	Duration uint `json:"duration,omitempty"`
	// Event that triggered the build
	EventType string `json:"event_type,omitempty"`
	// State of the previous build (useful to see if state changed)
	PreviousState string `json:"previous_state,omitempty"`
	// Title of the build's pull request
	PullRequestTitle string `json:"pull_request_title,omitempty"`
	// Number of the build's pull request
	PullRequestNumber uint `json:"pull_request_number,omitempty"`
	// When the build started
	StartedAt string `json:"started_at,omitempty"`
	// When the build finished
	FinishedAt string `json:"finished_at,omitempty"`
	// Whether or not the build is private
	Private bool `json:"private,omitempty"`
}

// buildResponse is only used to parse responses from Restart, Cancel
type buildResponse struct {
	Build MinimalBuild `json:"build,omitempty"`
}

const (
	// BuildStateCreated represents the build state `created`
	BuildStateCreated = "created"
	// BuildStateReceived represents the build state `received`
	BuildStateReceived = "received"
	// BuildStateStarted represents the build state `started`
	BuildStateStarted = "started"
	// BuildStatePassed represents the build state `passed`
	BuildStatePassed = "passed"
	// BuildStateFailed represents the build state `failed`
	BuildStateFailed = "failed"
	// BuildStateErrored represents the build state `errored`
	BuildStateErrored = "errored"
	// BuildStateCanceled represents the build state `canceled`
	BuildStateCanceled = "canceled"
)

const (
	// BuildEventTypePush represents the build event type `push`
	BuildEventTypePush = "push"
	// BuildEventTypePullRequest represents the build event type `pull_request`
	BuildEventTypePullRequest = "pull_request"
)

// Find fetches a build based on the provided build id
//
// Travis CI API docs: https://developer.travis-ci.com/resource/build#find
func (bs *BuildService) Find(ctx context.Context, id uint) (*Build, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/build/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest(http.MethodGet, u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var build Build
	resp, err := bs.client.Do(ctx, req, &build)
	if err != nil {
		return nil, resp, err
	}

	return &build, resp, err
}

// Cancel cancels a build based on the provided build id
//
// Travis CI API docs: https://developer.travis-ci.com/resource/build#cancel
func (bs *BuildService) Cancel(ctx context.Context, id uint) (*MinimalBuild, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/build/%d/cancel", id), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest(http.MethodPost, u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var response buildResponse
	resp, err := bs.client.Do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response.Build, resp, err
}

// Restart restarts a build based on the provided build id
//
// Travis CI API docs: https://developer.travis-ci.com/resource/build#restart
func (bs *BuildService) Restart(ctx context.Context, id uint) (*MinimalBuild, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/build/%d/restart", id), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest(http.MethodPost, u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var response buildResponse
	resp, err := bs.client.Do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response.Build, resp, err
}
