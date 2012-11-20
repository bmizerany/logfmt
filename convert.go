package logfmt

import (
	"reflect"
	"strings"
	"errors"
)

var ErrInvalidType = errors.New("logfmt: invalid type")

func assign(key string, x interface{}, tok *token) error {
	switch v := x.(type) {
	case map[string]string:
		v[key] = tok.string()
		return nil
	}

	sv := reflect.Indirect(reflect.ValueOf(x))
	if sv.Kind() != reflect.Struct {
		return ErrInvalidType
	}
	st := sv.Type()
	for i := 0; i < sv.NumField(); i++ {
		sf := st.Field(i)
		if strings.EqualFold(sf.Name, key) {
			return convertAssign(sv.Field(i), tok)
		}
	}
	return nil
}

// assumes dst.CanSet() == true
func convertAssign(dv reflect.Value, tok *token) error {
	if tok.isNull() {
		dv.Set(reflect.Zero(dv.Type()))
		return nil
	}

	for dv.Kind() == reflect.Ptr {
		dv.Set(reflect.New(dv.Type().Elem()))
		dv = reflect.Indirect(dv)
	}

	switch dv.Kind() {
	case reflect.String:
		dv.SetString(tok.string())
		return nil
	case reflect.Bool:
		dv.SetBool(tok.bool())
		return nil
	}

	switch {
	case reflect.Int <= dv.Kind() && dv.Kind() <= reflect.Int64:
		n, err := tok.int(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetInt(n)
		return nil
	case reflect.Uint <= dv.Kind() && dv.Kind() <= reflect.Uint64:
		n, err := tok.uint(dv.Type().Bits())
		if err != nil {
			return err
		}
		dv.SetUint(n)
		return nil
	}

	return nil
}
