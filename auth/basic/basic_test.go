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

package basicauth_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/auth"
	"github.com/mjpitz/myago/auth/basic"
	"github.com/mjpitz/myago/headers"
)

func TestBasic(t *testing.T) {
	testCases := []struct {
		Name          string
		Authorization string
		UserInfo      *auth.UserInfo
	}{
		{
			Name:          "empty",
			Authorization: "",
			UserInfo:      nil,
		},
		{
			Name:          "bearer",
			Authorization: "Bearer 52295778-63A0-4809-9B30-F3269D917114",
			UserInfo:      nil,
		},
		{
			Name:          "basic username",
			Authorization: "Basic dXNlcm5hbWU=",
			UserInfo:      nil,
		},
		{
			Name:          "basic admin:",
			Authorization: "Basic YWRtaW46",
			UserInfo:      nil,
		},
		{
			Name:          "basic username:",
			Authorization: "Basic dXNlcm5hbWU6",
			UserInfo:      nil,
		},
		{
			Name:          "basic username:password",
			Authorization: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			UserInfo: &auth.UserInfo{
				Subject:       "userID",
				Profile:       "username",
				Email:         "",
				EmailVerified: false,
				Groups:        []string{"group1", "group2"},
			},
		},
	}

	store, err := basicauth.OpenCSV(context.Background(), filepath.Join("testdata", "basic.csv"))
	require.NoError(t, err)

	handler := basicauth.Basic(store)

	for _, testCase := range testCases {
		t.Log(testCase.Name)

		header := make(headers.Header)
		if len(testCase.Authorization) > 0 {
			header.Set("authorization", testCase.Authorization)
		}

		ctx := headers.ToContext(context.Background(), header)
		ctx, err := handler(ctx)
		require.NoError(t, err)

		userInfo := auth.Extract(ctx)
		require.Equal(t, testCase.UserInfo, userInfo)
	}
}
