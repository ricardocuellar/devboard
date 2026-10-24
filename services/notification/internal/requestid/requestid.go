package requestid

import (
	"context"
	"crypto/rand"
	"regexp"
)

type contextKey struct{}

var validID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

func WithContext(ctx context.Context, id string) context.Context {
	if !validID.MatchString(id) {
		id = rand.Text()
	}

	return context.WithValue(ctx, contextKey{}, id)
}

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
