package tools

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type FieldMask struct {
	paths []string
}

func NewFieldMask(paths ...string) FieldMask {
	return FieldMask{
		paths: paths,
	}
}

func (m *FieldMask) Append(path string) {
	m.paths = append(m.paths, path)
}

func (m *FieldMask) Extract(x interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(x)
	v = reflect.Indirect(v)
	r := make([]interface{}, len(m.paths))

	for _, path := range m.paths {
		f := v.FieldByName(strings.ToTitle(path))
		if !f.IsValid() {
			return nil, errors.Errorf("cannot extract find path field (%s)", path)
		}
		r = append(r, f.Interface())
	}
	return r, nil
}
