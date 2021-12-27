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
	"os/user"
	"runtime"
	"strings"
)

// Must panics if Local fails to obtain the directory set.
func Must(appName string) DirectorySet {
	dirs, err := Local(appName)
	if err != nil {
		panic(err)
	}
	return dirs
}

// Local uses information about the local system to determine a common set of paths for use by the application.
func Local(appName string) (DirectorySet, error) {
	switch runtime.GOOS {
	case "windows":
		info, err := user.Current()
		if err != nil {
			return DirectorySet{}, err
		}

		return windows(info.HomeDir, appName), nil
	case "darwin":
		info, err := user.Current()
		if err != nil {
			return DirectorySet{}, err
		}

		return osx(info.HomeDir, appName), nil
	default:
		return linux(appName), nil
	}
}

func linux(appName string) DirectorySet {
	appName = strings.ToLower(appName)
	sep := "/"

	format := func(part string) string {
		return strings.Join([]string{"/var", part, appName}, sep)
	}

	return DirectorySet{
		CacheDir:      format("cache"),
		StateDir:      format("lib"),
		LocalStateDir: format("local"),
		LockDir:       format("locks"),
		LogDir:        format("logs"),
	}
}

func osx(home, appName string) DirectorySet {
	sep := "/"
	join := func(parts ...string) string {
		parts = append([]string{home, "Library", "Application Support", appName}, parts...)

		return strings.Join(parts, sep)
	}

	return DirectorySet{
		CacheDir:      join("Cache"),
		StateDir:      join("State"),
		LocalStateDir: join("Local"),
		LockDir:       join("Locks"),
		LogDir:        join("Logs"),
	}
}

func windows(home, appName string) DirectorySet {
	sep := "\\"
	join := func(parts ...string) string {
		parts = append([]string{home, "AppData", "Roaming", appName}, parts...)

		return strings.Join(parts, sep)
	}

	return DirectorySet{
		CacheDir:      join("Cache"),
		StateDir:      join("State"),
		LocalStateDir: join("Local"),
		LockDir:       join("Locks"),
		LogDir:        join("Logs"),
	}
}
