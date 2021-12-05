package wavl

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

func firstItem(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	return sl[0]
}

func spliceItems(sl []string, afterIdx int) []string {
	if afterIdx+1 > len(sl) {
		return []string{}
	}
	return sl[afterIdx+1:]
}

func itemAt(sl []string, idx int) string {
	if idx < 0 || idx >= len(sl) {
		return ""
	}
	return sl[idx]
}

func itemAtReverse(sl []string, ridx int) string {
	return itemAt(sl, len(sl)-ridx-1)
}

func scanItems(sl []string, vs ...interface{}) error {
	if len(vs) > len(sl) {
		return errors.Errorf("invalid amount (%d) of scan params (%d items). ", len(vs), len(sl))
	}
	for i, v := range vs {
		item := sl[i]
		err := scanItem(item, v)
		if err != nil {
			return errors.Wrapf(err, "scan-item %d", i)
		}
	}
	return nil
}

func scanItem(item string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("%T is not a pointer type", v)
	}
	elm := rv.Elem()
	if !elm.CanSet() {
		return errors.Errorf("cannot set %T", v)
	}

	switch elm.Kind() {
	case reflect.String:
		elm.SetString(item)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		n, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return err
		}
		elm.SetInt(n)
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		n, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			return err
		}
		elm.SetUint(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(item)
		if err != nil {
			return err
		}
		elm.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(item, 64)
		if err != nil {
			return err
		}
		elm.SetFloat(f)
	default:
		return errors.Errorf("cannot scan into value of type %T", v)
	}

	return nil
}
