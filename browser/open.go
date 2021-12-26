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
