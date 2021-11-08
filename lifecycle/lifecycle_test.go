package lifecycle_test

import (
	"context"
	"sync"
	"testing"

	"github.com/mjpitz/myago/lifecycle"
)

func TestLifeCycle(t *testing.T) {
	ctx := context.Background()

	lc := &lifecycle.LifeCycle{}
	ctx = lc.Setup(ctx)

	defer lc.Shutdown(ctx)

	wg := sync.WaitGroup{}

	{
		wg.Add(3)

		for i := 0; i < 3; i++ {
			lc.Defer(func(ctx context.Context) {
				wg.Done()
			})
		}

		lc.Resolve(context.Background())
		wg.Wait()
	}

	{
		wg.Add(3)

		for i := 0; i < 3; i++ {
			lc.Defer(func(ctx context.Context) {
				wg.Done()
			})
		}

		lc.Resolve(context.Background())
		wg.Wait()
	}
}
