package record

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

func (d *Decoder) Decode(value interface{}) error {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Ptr:
		if v.IsNil() || v.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("record: invalid pointer") // TODO: move to type const
		}
		return d.decodeStruct(v.Elem(), v.Elem().Type())
	case reflect.Struct:
		return d.decodeStruct(v, t)
	}
	return fmt.Errorf("record: invalid value type: %s", t)
}

func (d *Decoder) decodeStruct(v reflect.Value, t reflect.Type) error {
	var (
		l     string        // line
		token string        // next token
		start int           // start
		tag   *tag          // struct tag with metadata
		fval  reflect.Value // field to set
		err   error
	)

	if l, err = d.readLine(); err != nil {
		return err
	}

	for i := 0; i < t.NumField() && start < len(l); i++ {
		f := t.Field(i)
		if tag = parseTags(f); tag.skip {
			continue
		}

		if token, start, err = nextToken(l, start, tag); err != nil {
			return err
		}

		if fval = v.Field(i); !fval.CanSet() {
			log.Printf("can't set %v", f)
			continue
		}

		switch f.Type.Kind() {
		case reflect.String:
			fval.SetString(strings.TrimRight(token, " \t\n\r"))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(token, 10, 64); err != nil {
				return err
			} else {
				fval.SetInt(intVal)
			}
		case reflect.Struct:
			if f.Type.ConvertibleTo(dateType) {
				// We need to parse, reformat into the MarshallText, then unmarshal in t again
				if timeVal, err := time.Parse(DateFormat, token); err != nil {
					return err
				} else {
					fval.Set(reflect.ValueOf(timeVal))
				}
			} else {
				// Attempt to deep-decode nested structs
				return d.decodeStruct(fval, f.Type)
			}
		default:
			return fmt.Errorf("record: unsupported type: %v", f)
		}
	}

	return nil
}

func nextToken(l string, start int, t *tag) (string, int, error) {
	if t == nil {
		return "", start, fmt.Errorf("record: unexpected nil tag")
	}
	end := start + t.size
	if t.size == 0 {
		end = len(l) - 1
	}
	if end > len(l) {
		return "", end, fmt.Errorf("record: end of line")
	}
	return string([]rune(l)[start:end]), end, nil
}

func (d *Decoder) readLine() (string, error) {
	return d.r.ReadString('\n')
}

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
