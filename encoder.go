package record

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Marshal takes a struct value or pointer
// and returns the encoded bytes and a nil error,
// or a nil byte slice and the encoding error.
func Marshal(src interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(src); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Encoder is responsible for encoding record values
// from struct fields.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns an initialized encoder
// that writes the encoded bytes into the specified io.Writer w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode takes an struct value or pointer
// and encode its fields into the encoder writer.
// Returns any encoding errors, if any.
// When returning an error, src may have been partially written.
func (e *Encoder) Encode(src interface{}) error {
	v := reflect.ValueOf(src)
	t := reflect.TypeOf(src)

	switch t.Kind() {
	case reflect.Ptr:
		if v.IsNil() || v.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("record: invalid pointer")
		}
		return e.encodeStruct(v.Elem(), v.Elem().Type())
	case reflect.Struct:
		return e.encodeStruct(v, t)
	}
	return fmt.Errorf("record: invalid value type: %s", t)
}

// encodeStruct encodes the content of struct s into w.
func (e *Encoder) encodeStruct(val reflect.Value, sType reflect.Type) error {
	for i := 0; i < sType.NumField(); i++ {
		f := sType.Field(i)
		tag := parseTags(f)
		if tag.skip {
			continue
		}
		switch f.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmtSpec := "%d"
			if !tag.noPadding {
				fmtSpec = fmt.Sprintf("%%0%dd", tag.size)
			}
			fmt.Fprintf(e.w, fmtSpec, val.Field(i).Int())
		case reflect.String:
			fmtSpec := "%s"
			if !tag.noPadding {
				fmtSpec = fmt.Sprintf("%% %ds", tag.size)
			}
			str := val.Field(i).String()
			if tag.upper {
				str = strings.ToUpper(str)
			}
			if len(str) > tag.size && tag.size > 0 {
				str = string([]rune(str)[:tag.size])
			}
			fmt.Fprintf(e.w, fmtSpec, str)
		}
	}
	return nil
}

// tag holds configuration options for encoding fields.
type tag struct {
	size      int
	noPadding bool
	upper     bool
	skip      bool
}

// parseTags takes a field and returns the parsed tag.
// The default tag values are returned, when no tag is specified.
func parseTags(f reflect.StructField) *tag {
	t := &tag{}
	tagVal := f.Tag.Get("record")
	if tagVal != "" {
		elem := strings.Split(f.Tag.Get("record"), ",")
		if size, err := strconv.Atoi(elem[0]); err == nil {
			t.size = size
			elem = elem[1:]
		}
		for _, e := range elem {
			switch e {
			case "nopad", "nopadding":
				t.noPadding = true
			case "upper":
				t.upper = true
			case "-":
				t.skip = true
			}
		}
	}
	return t
}
