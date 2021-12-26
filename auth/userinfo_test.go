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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/auth"
)

type Extensions struct {
	Cubbyhole string `json:"cubbyhole"`
}

type CustomUserInfo struct {
	auth.UserInfo
	Extensions
}

const expected = `{
  "sub": "0x01234",
  "profile": "jane doe",
  "email": "jane@example.com",
  "email_verified": true,
  "groups": [
    "group1",
    "group2"
  ],
  "cubbyhole": "blah.blah.blah"
}`

func TestUserInfo(t *testing.T) {
	customInfo := &CustomUserInfo{
		UserInfo: auth.UserInfo{
			Subject:       "0x01234",
			Profile:       "jane doe",
			Email:         "jane@example.com",
			EmailVerified: true,
			Groups:        []string{"group1", "group2"},
		},
		Extensions: Extensions{
			Cubbyhole: "blah.blah.blah",
		},
	}

	body, err := json.MarshalIndent(customInfo, "", "  ")
	require.NoError(t, err)
	require.Equal(t, expected, string(body))

	base := auth.UserInfo{}
	ext := Extensions{}

	err = json.Unmarshal(body, &base)
	require.NoError(t, err)

	require.Equal(t, customInfo.UserInfo.Subject, base.Subject)
	require.Equal(t, customInfo.UserInfo.Profile, base.Profile)
	require.Equal(t, customInfo.UserInfo.Email, base.Email)
	require.Equal(t, customInfo.UserInfo.EmailVerified, base.EmailVerified)
	require.Equal(t, customInfo.UserInfo.Groups, base.Groups)

	err = base.Claims(&ext)
	require.NoError(t, err)

	require.Equal(t, customInfo.Extensions.Cubbyhole, ext.Cubbyhole)
}
