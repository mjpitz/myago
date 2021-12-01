package ulid

import (
	"context"
	"os"
	"strconv"

	"github.com/mjpitz/myago"
)

var contextKey = myago.ContextKey("ulid.generator")
var systemGenerator *Generator

func init() {
	skew := byte(1)
	if skewEnv := os.Getenv("MYAGO_ULID_SKEW"); skewEnv != "" {
		s, err := strconv.Atoi(skewEnv)
		if err != nil {
			// not a fan
			panic(err)
		}
		skew = byte(s)
	}

	systemGenerator = NewGenerator(skew, RandomFill)
}

// Extract is used to obtain the generator from a context. If none is present, the system generator is used.
func Extract(ctx context.Context) *Generator {
	val := ctx.Value(contextKey)
	if val == nil {
		return systemGenerator
	}

	return val.(*Generator)
}

// ToContext appends the provided generator to the provided context.
func ToContext(ctx context.Context, generator *Generator) context.Context {
	return context.WithValue(ctx, contextKey, generator)
}
