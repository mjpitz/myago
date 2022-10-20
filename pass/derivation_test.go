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

package pass_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/pass"
)

func TestDerivation(t *testing.T) {
	t.Parallel()

	password := []byte("badadmin")
	user := "admin"

	testCases := []struct {
		site     string
		scope    pass.Scope
		counter  uint32
		template pass.TemplateClass
		password string
	}{
		{"scope.varies", pass.Authentication, 0, pass.MaximumSecurity, "Yc(BMzvIdhLt*JQhPZ3~"},
		{"scope.varies", pass.Identification, 0, pass.MaximumSecurity, "B2%%OeMJuiruYZ$un34s"},
		{"scope.varies", pass.Recovery, 0, pass.MaximumSecurity, "jf3WgS2PMu$35fbNbG1^"},
		{"counter.varies", pass.Authentication, 3, pass.MaximumSecurity, "a(I3rJe&qDBVS^uFF@3!"},
		{"counter.varies", pass.Authentication, 5, pass.MaximumSecurity, "Rm4HJRqHMA$fxUNnoK8#"},
		{"counter.varies", pass.Authentication, math.MaxInt32, pass.MaximumSecurity, "S1&Pcjik6*bGTk!U$#*V"},
		{"template.varies", pass.Authentication, 0, pass.MaximumSecurity, "edqYl3g7pDuj3lf9t08="},
		{"template.varies", pass.Authentication, 0, pass.Long, "HansJaduXutc2!"},
		{"template.varies", pass.Authentication, 0, pass.Medium, "HanSer9+"},
		{"template.varies", pass.Authentication, 0, pass.Short, "Han3"},
		{"template.varies", pass.Authentication, 0, pass.Basic, "eE83iOO6"},
		{"template.varies", pass.Authentication, 0, pass.PIN, "3783"},
	}

	for _, testCase := range testCases {
		identity, err := pass.Identity(testCase.scope, password, user)
		require.NoError(t, err)

		siteKey := pass.SiteKey(testCase.scope, identity, testCase.site, testCase.counter)
		password := pass.SitePassword(siteKey, testCase.template)

		require.Equal(t, testCase.password, string(password))
	}
}
