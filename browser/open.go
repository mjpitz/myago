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

package browser

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// Open attempts to open the provided url in a browser.
func Open(ctx context.Context, url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.CommandContext(ctx, "xdg-open", url).Start()
	case "windows":
		return exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.CommandContext(ctx, "open", url).Start()
	}

	return fmt.Errorf("unsupported platform")
}
