package config

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func TestConfigLoad(t *testing.T) {
	filepath := "./testdata/config.yml"

	conf, err := NewFromFile(filepath)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	t.Logf("%s", spew.Sdump(conf))
}
