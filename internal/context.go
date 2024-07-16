package internal

import "context"

// InitIfNilContext initializes a context if the provided context is nil.
// It takes a context.Context as input and returns a context.Context.
// If the input context is nil (determined by the IsNil function),
// it returns a new background context. Otherwise, it returns the input context.
func InitIfNilContext(ctx context.Context) context.Context {
	if IsNil(ctx) {
		return context.Background()
	}
	return ctx
}
