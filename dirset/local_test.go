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

package dirset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinux(t *testing.T) {
	t.Parallel()

	directorySet := linux("ExampleApp")

	require.Equal(t, "/var/cache/exampleapp", directorySet.CacheDir)
	require.Equal(t, "/var/lib/exampleapp", directorySet.StateDir)
	require.Equal(t, "/var/local/exampleapp", directorySet.LocalStateDir)
	require.Equal(t, "/var/locks/exampleapp", directorySet.LockDir)
	require.Equal(t, "/var/logs/exampleapp", directorySet.LogDir)
}

func TestOSX(t *testing.T) {
	t.Parallel()

	directorySet := osx("/Users/myago", "ExampleApp")

	require.Equal(t, "/Users/myago/Library/Application Support/ExampleApp/Cache", directorySet.CacheDir)
	require.Equal(t, "/Users/myago/Library/Application Support/ExampleApp/State", directorySet.StateDir)
	require.Equal(t, "/Users/myago/Library/Application Support/ExampleApp/Local", directorySet.LocalStateDir)
	require.Equal(t, "/Users/myago/Library/Application Support/ExampleApp/Locks", directorySet.LockDir)
	require.Equal(t, "/Users/myago/Library/Application Support/ExampleApp/Logs", directorySet.LogDir)
}

func TestWindows(t *testing.T) {
	t.Parallel()

	directorySet := windows("C:\\Users\\myago", "ExampleApp")

	require.Equal(t, "C:\\Users\\myago\\AppData\\Roaming\\ExampleApp\\Cache", directorySet.CacheDir)
	require.Equal(t, "C:\\Users\\myago\\AppData\\Roaming\\ExampleApp\\State", directorySet.StateDir)
	require.Equal(t, "C:\\Users\\myago\\AppData\\Roaming\\ExampleApp\\Local", directorySet.LocalStateDir)
	require.Equal(t, "C:\\Users\\myago\\AppData\\Roaming\\ExampleApp\\Locks", directorySet.LockDir)
	require.Equal(t, "C:\\Users\\myago\\AppData\\Roaming\\ExampleApp\\Logs", directorySet.LogDir)
}
