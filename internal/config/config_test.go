package config

import (
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func TestConfigLoad(t *testing.T) {
	config := NewDefault()

	filepath := "./testdata/config.yml"

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(errors.Wrapf(err, "could not read file '%s'", filepath))
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		t.Fatal(errors.Wrapf(err, "could not unmarshal configuration"))
	}

	// t.Logf("%s", spew.Sdump(config))
}
