package config

import (
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var reVar = regexp.MustCompile(`^\${(\w+)}$`)

type InterpolatedString string

func (is *InterpolatedString) UnmarshalYAML(value *yaml.Node) error {
	var str string

	if err := value.Decode(&str); err != nil {
		return errors.WithStack(err)
	}

	if match := reVar.FindStringSubmatch(str); len(match) > 0 {
		*is = InterpolatedString(os.Getenv(match[1]))
	} else {
		*is = InterpolatedString(str)
	}

	return nil
}

type InterpolatedInt int

func (ii *InterpolatedInt) UnmarshalYAML(value *yaml.Node) error {
	var str string

	if err := value.Decode(&str); err != nil {
		return errors.Wrapf(err, "could not decode value '%v' (line '%d') into string", value.Value, value.Line)
	}

	if match := reVar.FindStringSubmatch(str); len(match) > 0 {
		str = os.Getenv(match[1])
	}

	intVal, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return errors.Wrapf(err, "could not parse int '%v',  line '%d'", str, value.Line)
	}

	*ii = InterpolatedInt(int(intVal))

	return nil
}

type InterpolatedBool bool

func (ib *InterpolatedBool) UnmarshalYAML(value *yaml.Node) error {
	var str string

	if err := value.Decode(&str); err != nil {
		return errors.Wrapf(err, "could not decode value '%v' (line '%d') into string", value.Value, value.Line)
	}

	if match := reVar.FindStringSubmatch(str); len(match) > 0 {
		str = os.Getenv(match[1])
	}

	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		return errors.Wrapf(err, "could not parse bool '%v',  line '%d'", str, value.Line)
	}

	*ib = InterpolatedBool(boolVal)

	return nil
}

type InterpolatedMap map[string]interface{}

func (im *InterpolatedMap) UnmarshalYAML(value *yaml.Node) error {
	var data map[string]interface{}

	if err := value.Decode(&data); err != nil {
		return errors.Wrapf(err, "could not decode value '%v' (line '%d') into map", value.Value, value.Line)
	}

	for key, value := range data {
		strVal, ok := value.(string)
		if !ok {
			continue
		}

		if match := reVar.FindStringSubmatch(strVal); len(match) > 0 {
			strVal = os.Getenv(match[1])
		}

		data[key] = strVal
	}

	*im = data

	return nil
}
