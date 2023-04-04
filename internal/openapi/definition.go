package openapi

import (
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
)

func Definition(t reflect.Type) *spec.Schema {
	switch t.Kind() {
	case reflect.String:
		return spec.StringProperty()
	case reflect.Int64, reflect.Int:
		return spec.Int64Property()
	case reflect.Int32:
		return spec.Int32Property()
	case reflect.Int16:
		return spec.Int16Property()
	case reflect.Int8:
		return spec.Int8Property()
	case reflect.Float64:
		return spec.Float64Property()
	case reflect.Float32:
		return spec.Float32Property()
	case reflect.Pointer:
		return Definition(t.Elem())
	case reflect.Struct:
		s := &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Properties: spec.SchemaProperties{},
			},
		}

		for i := 0; i < t.NumField(); i++ {
			var (
				sf  = t.Field(i)
				key = sf.Name
			)

			if sf.IsExported() {
				k := strings.SplitN(sf.Tag.Get("json"), ",", 2)[0]
				if k == "-" {
					continue
				} else if k != "" {
					key = k
				}

				s.SchemaProps.Properties[key] = *Definition(sf.Type)
			}
		}

		return spec.MapProperty(s)
	case reflect.Array, reflect.Slice:
		return spec.ArrayProperty(Definition(t.Elem()))
	}

	return nil
}
