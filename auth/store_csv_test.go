// Copyright (C) 2021 Mya Pitzeruse
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

package auth_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/auth"
)

type testUser struct {
	Request  auth.LookupRequest
	Response auth.LookupResponse
	Error    string
}

func TestCSVBasicStore(t *testing.T) {
	ctx := context.Background()

	store, err := auth.OpenCSV(ctx, filepath.Join("testdata", "basic.csv"))
	require.NoError(t, err)

	testUsers := []testUser{
		{
			Request: auth.LookupRequest{
				User: "username",
			},
			Response: auth.LookupResponse{
				Password: "password",
				User:     "username",
				UserID:   "userID",
				Groups:   []string{"group1", "group2"},
			},
		},
		{
			Request: auth.LookupRequest{
				Token: "invalid",
			},
			Error: "not found",
		},
	}

	for _, testUser := range testUsers {
		resp, err := store.Lookup(testUser.Request)
		if len(testUser.Error) > 0 {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, testUser.Response, resp)
		}
	}
}

func TestCSVTokenStore(t *testing.T) {
	ctx := context.Background()

	store, err := auth.OpenCSV(ctx, filepath.Join("testdata", "token.csv"))
	require.NoError(t, err)

	testUsers := []testUser{
		{
			Request: auth.LookupRequest{
				Token: "token",
			},
			Response: auth.LookupResponse{
				User:   "username",
				Token:  "token",
				UserID: "userID",
				Groups: []string{"group1"},
			},
		},
		{
			Request: auth.LookupRequest{
				Token: "invalid",
			},
			Error: "not found",
		},
	}

	for _, testUser := range testUsers {
		resp, err := store.Lookup(testUser.Request)
		if len(testUser.Error) > 0 {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, testUser.Response, resp)
		}
	}
}
