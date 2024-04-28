package xlsx

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func Open(r io.Reader) (*excelize.File, error) {
	return excelize.OpenReader(r)
}

func GetAxis(x, y int) string {
	return Num2AZ(x) + strconv.Itoa(y)
}

func Num2AZ(num int) string {
	var (
		str  string
		k    int
		temp []int
	)
	slices := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	if num > 26 {
		for {
			k = num % 26
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			num = (num - k) / 26
			if num <= 26 {
				temp = append(temp, num)
				break
			}
		}
	} else {
		return slices[num]
	}
	for _, value := range temp {
		str = slices[value] + str
	}
	return str
}

func String2Value(s string, rv reflect.Value) error {
	var (
		val interface{}
		err error
	)
	if s == "" {
		val = reflect.Zero(rv.Type()).Interface()
	} else {
		switch rv.Kind() {
		case reflect.Bool:
			switch strings.ToLower(s) {
			case "false", "f", "0":
				val = false
			case "true", "t", "1":
				val = true
			default:
				return fmt.Errorf("parse %s to bool error", s)
			}
		case reflect.Int:
			val, err = strconv.Atoi(s)
		case reflect.Int8:
			t, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return err
			}
			val = int8(t)
		case reflect.Int16:
			t, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return err
			}
			val = int16(t)
		case reflect.Int32:
			t, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return err
			}
			val = int32(t)
		case reflect.Int64:
			val, err = strconv.ParseInt(s, 10, 64)
		case reflect.Uint:
			t, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return err
			}
			val = uint(t)
		case reflect.Uint8:
			t, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return err
			}
			val = uint8(t)
		case reflect.Uint16:
			t, err := strconv.ParseUint(s, 10, 16)
			if err != nil {
				return err
			}
			val = uint16(t)
		case reflect.Uint32:
			t, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return err
			}
			val = uint32(t)
		case reflect.Uint64:
			val, err = strconv.ParseUint(s, 10, 64)
		case reflect.Float32:
			t, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return err
			}
			val = float32(t)
		case reflect.Float64:
			val, err = strconv.ParseFloat(s, 64)
		case reflect.String:
			val = s
		default:
			return errors.New("not Supported " + rv.Kind().String())
		}
	}
	if err != nil {
		return err
	}
	rv.Set(reflect.ValueOf(val))
	return nil
}

func ZeroValue(kind reflect.Kind) (interface{}, error) {
	var v interface{}
	switch kind {
	case reflect.Bool:
		v = false
	case reflect.Int:
		v = int(0)
	case reflect.Int8:
		v = int8(0)
	case reflect.Int16:
		v = int16(0)
	case reflect.Int32:
		v = int32(0)
	case reflect.Int64:
		v = int64(0)
	case reflect.Uint:
		v = uint(0)
	case reflect.Uint8:
		v = uint8(0)
	case reflect.Uint16:
		v = uint16(0)
	case reflect.Uint32:
		v = uint32(0)
	case reflect.Uint64:
		v = uint64(0)
	case reflect.Float32:
		v = float32(0)
	case reflect.Float64:
		v = float64(0)
	case reflect.String:
		v = ""
	default:
		return nil, errors.New("not Supported " + kind.String())
	}
	return v, nil
}

func GetSheetName(v interface{}) string {
	return getSheetName(reflect.TypeOf(v))
}

func getSheetName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}
	if s, ok := reflect.New(t).Interface().(SheetName); ok {
		return s.SheetName()
	} else {
		return t.Name()
	}
}
