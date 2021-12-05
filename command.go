package wavx

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Command map[string]string

func (c Command) String() string {
	var sl []string
	for k, v := range c {
		sl = append(sl, fmt.Sprintf("%s:%s", k, v))
	}
	sort.Strings(sl)
	return strings.Join(sl, ", ")
}

type ChangeMode int

const (
	ChangeAbs ChangeMode = iota
	ChangeInc
	ChangeDec
)

type ChangeFloat struct {
	Mode  ChangeMode
	Value float64
}

func (cf ChangeFloat) Applied(v float64) float64 {
	switch cf.Mode {
	case ChangeInc:
		return v + cf.Value
	case ChangeDec:
		return v - cf.Value
	default:
		return cf.Value
	}
}

func ParseChangeFloat(s string) (ChangeFloat, error) {
	mode := ChangeAbs
	if strings.HasPrefix(s, "+") {
		mode = ChangeInc
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		mode = ChangeDec
		s = s[1:]
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ChangeFloat{}, err
	}
	return ChangeFloat{
		Mode:  mode,
		Value: f,
	}, nil
}

type ChangeInt struct {
	Mode  ChangeMode
	Value int64
}

func (ci ChangeInt) Applied(v int64) int64 {
	switch ci.Mode {
	case ChangeInc:
		return v + ci.Value
	case ChangeDec:
		return v - ci.Value
	default:
		return ci.Value
	}
}

func ParseChangeInt(s string) (ChangeInt, error) {
	mode := ChangeAbs
	if strings.HasPrefix(s, "+") {
		mode = ChangeInc
		s = s[1:]
	} else if strings.HasPrefix(s, "-") {
		mode = ChangeDec
		s = s[1:]
	}

	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ChangeInt{}, err
	}
	return ChangeInt{
		Mode:  mode,
		Value: f,
	}, nil
}

func ApplyCommand(cmd Command, data interface{}) error {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr {
		return errors.Errorf("%T is not a pointer type", data)
	}
	elm := rv.Elem()
	if elm.Kind() != reflect.Struct {
		return errors.Errorf("%T is not a struct", data)
	}

	for k, s := range cmd {
		fv := elm.FieldByNameFunc(func(name string) bool {
			return strings.ToLower(name) == k
		})
		if !fv.CanSet() {
			return errors.Errorf("cannot set %q", k)
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(s)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ci, err := ParseChangeInt(s)
			if err != nil {
				return err
			}
			fv.SetInt(ci.Applied(fv.Int()))
		case
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ci, err := ParseChangeInt(s)
			if err != nil {
				return err
			}
			fv.SetUint(uint64(ci.Applied(fv.Int())))
		case reflect.Bool:
			b, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			fv.SetBool(b)
		case reflect.Float32, reflect.Float64:
			cf, err := ParseChangeFloat(s)
			if err != nil {
				return err
			}
			fv.SetFloat(cf.Applied(fv.Float()))
		}
	}

	return nil
}
