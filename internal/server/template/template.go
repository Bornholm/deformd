package template

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"path"

	"github.com/Masterminds/sprig/v3"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	layouts map[string]*template.Template
	blocks  map[string]string
)

func Load(files embed.FS, baseDir string) error {
	if blocks == nil {
		blocks = make(map[string]string)
	}

	blockFiles, err := fs.ReadDir(files, path.Join(baseDir, "blocks"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range blockFiles {
		templateData, err := fs.ReadFile(files, path.Join(baseDir, "blocks/"+f.Name()))
		if err != nil {
			return errors.WithStack(err)
		}

		blocks[f.Name()] = string(templateData)
	}

	layoutFiles, err := fs.ReadDir(files, path.Join(baseDir, "layouts"))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, f := range layoutFiles {
		templateData, err := fs.ReadFile(files, path.Join(baseDir, "layouts/"+f.Name()))
		if err != nil {
			return errors.WithStack(err)
		}

		if err := loadLayout(f.Name(), string(templateData)); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func loadLayout(name string, rawTemplate string) error {
	if layouts == nil {
		layouts = make(map[string]*template.Template)
	}

	tmpl := template.New(name)
	funcMap := mergeHelpers(
		sprig.FuncMap(),
		customHelpers(tmpl),
	)

	tmpl.Funcs(funcMap)

	for blockName, b := range blocks {
		if _, err := tmpl.Parse(b); err != nil {
			return errors.Wrapf(err, "could not parse template block '%s'", blockName)
		}
	}

	tmpl, err := tmpl.Parse(rawTemplate)
	if err != nil {
		return errors.Wrapf(err, "could not parse template '%s'", name)
	}

	layouts[name] = tmpl

	return nil
}

func Exec(name string, w io.Writer, data interface{}) error {
	tmpl, exists := layouts[name]
	if !exists {
		return errors.Errorf("could not find template '%s'", name)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func mergeHelpers(helpers ...template.FuncMap) template.FuncMap {
	merged := template.FuncMap{}

	for _, help := range helpers {
		for name, fn := range help {
			merged[name] = fn
		}
	}

	return merged
}

type FormItemData struct {
	Parent         *FormItemData
	Schema         *jsonschema.Schema
	Property       string
	Error          *jsonschema.ValidationError
	Values         interface{}
	Defaults       interface{}
	SuccessMessage string
}

func customHelpers(tpl *template.Template) template.FuncMap {
	return template.FuncMap{
		"map": func(args ...interface{}) (map[string]interface{}, error) {
			data := make(map[string]interface{})

			var (
				key string
				ok  bool
			)

			for index, val := range args {
				if index%2 == 0 {
					key, ok = val.(string)
					if !ok {
						return nil, errors.Errorf("argument #%d must be a string to be used as key", index)
					}
				} else {
					data[key] = val
				}
			}

			return data, nil
		},
		"dump": func(data interface{}) string {
			spew.Dump(data)

			return ""
		},
		"include": func(name string, data interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer([]byte{})

			if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
				return "", errors.WithStack(err)
			}

			return template.HTML(buf.String()), nil
		},
	}
}
