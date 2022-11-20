package handler_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/pkg/errors"
)

func TestHandlerTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}

	t.Parallel()

	script := `while(true) {}`

	handler := handler.New(script, handler.WithMaxDuration(5*time.Second))

	ctx := context.Background()

	err := handler.Process(ctx, url.Values{})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(errors.Wrap(err, "err should be deadline exceeded"))
	}
}
