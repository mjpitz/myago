// Copyright (C) 2022 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	baseURL string
	client  *http.Client
}

func (api *API) Do(ctx context.Context, method, url string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, method, api.baseURL+url, body)
	if err != nil {
		return nil, err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	switch {
	case resp.StatusCode > 300:
		return nil, fmt.Errorf(resp.Status)
	}

	return resp.Body, nil
}

func (api *API) Jobs() *Jobs {
	return &Jobs{
		api: api,
	}
}

func (api *API) Pipelines() *Pipelines {
	return &Pipelines{
		api: api,
	}
}

type Jobs struct {
	api *API
}

func (jobs *Jobs) List(ctx context.Context) (*ListJobsResponse, error) {
	resp, err := jobs.api.Do(ctx, http.MethodGet, "/api/json", nil)
	if err != nil {
		return nil, err
	}

	response := &ListJobsResponse{}

	err = json.NewDecoder(resp).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (jobs *Jobs) Get(ctx context.Context, project string) (*GetJobResponse, error) {
	url := fmt.Sprintf("/job/%s/api/json", project)

	resp, err := jobs.api.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response := &GetJobResponse{}

	err = json.NewDecoder(resp).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type Pipelines struct {
	api *API
}

func (pipelines *Pipelines) Get(ctx context.Context, project string, build int) (*GetPipelineResponse, error) {
	url := fmt.Sprintf("/job/%s/%d/wfapi/describe", project, build)

	resp, err := pipelines.api.Do(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response := &GetPipelineResponse{}

	err = json.NewDecoder(resp).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
