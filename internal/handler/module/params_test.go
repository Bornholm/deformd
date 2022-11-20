package module_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/Bornholm/deformd/internal/handler/module"
	"github.com/pkg/errors"
)

func TestParamsModule(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Parallel()

	ctx := context.Background()

	script, err := loadTestDataScript("params.js")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	handler := handler.New(
		script,
		handler.WithModules(
			module.ParamsModuleFactory(map[string]interface{}{
				"foo": map[string]int{
					"bar": 1,
				},
			}),
		),
		handler.WithMaxDuration(5*time.Second),
	)

	form := url.Values{}

	if err := handler.Process(ctx, form); err != nil {
		t.Error(errors.WithStack(err))
	}
}
