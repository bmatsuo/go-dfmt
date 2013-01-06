package deepfmt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"unsafe"
)

var (
	newlineBytes     = []byte("\n")
	indentBytes      = []byte("\t")
	leftcurlyBytes   = []byte("{")
	rightcurlyBytes  = []byte("}")
	leftsquareBytes  = []byte("[")
	rightsquareBytes = []byte("]")
	leftparenBytes   = []byte("(")
	rightparenBytes  = []byte(")")
	colonBytes       = []byte(":")
	commaBytes       = []byte(",")
	spaceBytes       = []byte(" ")
	ptrBytes         = []byte("&")
	nilBytes         = []byte("nil")
	nilangleBytes    = []byte("<nil>")
	mapBytes         = []byte("map")
)

func writeNewline(s fmt.State)               { s.Write(newlineBytes) }
func writeIndent(s fmt.State)                { s.Write(indentBytes) }
func writeLeftcurly(s fmt.State)             { s.Write(leftcurlyBytes) }
func writeRightcurly(s fmt.State)            { s.Write(rightcurlyBytes) }
func writeLeftsquare(s fmt.State)            { s.Write(leftsquareBytes) }
func writeRightsquare(s fmt.State)           { s.Write(rightsquareBytes) }
func writeLeftparen(s fmt.State)             { s.Write(leftparenBytes) }
func writeRightparen(s fmt.State)            { s.Write(rightparenBytes) }
func writeColon(s fmt.State)                 { s.Write(colonBytes) }
func writeComma(s fmt.State)                 { s.Write(commaBytes) }
func writeSpace(s fmt.State)                 { s.Write(spaceBytes) }
func writePtr(s fmt.State)                   { s.Write(ptrBytes) }
func writeNil(s fmt.State)                   { s.Write(nilBytes) }
func writeNilangle(s fmt.State)              { s.Write(nilangleBytes) }
func writeMap(s fmt.State)                   { s.Write(mapBytes) }
func writeType(s fmt.State, v reflect.Value) { s.Write([]byte(v.Type().String())) }
func writeFullIndent(s fmt.State, n int) {
	for i := 0; i < n; i++ {
		writeIndent(s)
	}
}

type formatter struct {
	depth   int
	verbose bool
	pretty  bool
	v       interface{}
}

func NewFormatter(v interface{}) fmt.Formatter {
	return &formatter{v: v}
}

func (f *formatter) Format(s fmt.State, c rune) {
	if c != 'v' || !s.Flag('+') {
		newfmt := reconstructFlags(s, c)
		fmt.Fprintf(s, newfmt, f.v)
		return
	}
	if f.v == nil {
		writeNilangle(s)
		return
	}

	f.verbose = s.Flag('#')
	f.pretty = s.Flag(' ')
	f.format(s, c, reflect.ValueOf(f.v))
}

func (f *formatter) format(s fmt.State, c rune, val reflect.Value) {
	if !f.verbose {
		switch v := val.Interface(); v.(type) {
		case fmt.Stringer:
			fmt.Fprint(s, v.(fmt.Stringer).String())
			return
		case error:
			fmt.Fprint(s, v.(error).Error())
			return
		}
	} else if gs, ok := val.Interface().(fmt.GoStringer); ok {
		fmt.Fprint(s, gs.GoString())
		return
	}

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			if f.verbose {
				writeLeftparen(s)
				writeType(s, val)
				writeRightparen(s)
				writeLeftparen(s)
				writeNil(s)
				writeRightparen(s)
			} else {
				writeNilangle(s)
			}
			return
		}
		writePtr(s)
		f.format(s, c, reflect.Indirect(val))
	case reflect.Array, reflect.Slice:
		f.depth++
		if f.verbose {
			writeType(s, val)
			writeLeftcurly(s)
			if f.pretty {
				writeNewline(s)
				writeFullIndent(s, f.depth)
			}
		} else {
			writeLeftsquare(s)
		}
		for i, n := 0, val.Len(); i < n; i++ {
			if i > 0 {
				f.sep(s)
			}
			f.format(s, c, val.Index(i))
		}
		f.depth--
		if f.verbose {
			if f.pretty {
				writeComma(s)
				writeNewline(s)
				writeFullIndent(s, f.depth)
			}
			writeRightcurly(s)
		} else {
			writeRightsquare(s)
		}
	case reflect.Map:
		if val.IsNil() {
			writeNil(s)
			return
		}
		f.depth++
		if f.verbose {
			writeType(s, val)
			writeLeftcurly(s)
			if f.pretty {
				writeNewline(s)
				writeFullIndent(s, f.depth)
			}
		} else {
			writeMap(s)
			writeLeftsquare(s)
		}
		for i, key := range val.MapKeys() {
			if i > 0 {
				f.sep(s)
			}
			f.format(s, c, key)
			writeColon(s)
			f.format(s, c, val.MapIndex(key))
		}
		f.depth--
		if f.verbose {
			if f.pretty {
				writeComma(s)
				writeNewline(s)
				writeFullIndent(s, f.depth)
			}
			writeRightcurly(s)
		} else {
			writeRightsquare(s)
		}
	case reflect.Struct:
		f.depth++
		if f.verbose {
			writeType(s, val)
			writeLeftcurly(s)
			if f.pretty {
				writeNewline(s)
				writeFullIndent(s, f.depth)
			}
		} else {
			writeLeftcurly(s)
		}
		typ := val.Type()
		for i, n := 0, val.NumField(); i < n; i++ {
			field := typ.Field(i)
			if i > 0 {
				f.sep(s)
			}
			if f.verbose {
				s.Write([]byte(field.Name))
				writeColon(s)
			}
			if field.PkgPath == "" {
				f.format(s, c, val.Field(i))
			} else {
				field := typ.Field(i)
				var valptr reflect.Value
				if val.CanAddr() {
					valptr = val.Addr()
				} else {
					valptr = reflect.New(typ)
					reflect.Indirect(valptr).Set(val)
				}
				fieldp := valptr.Pointer() + field.Offset
				fieldptr := reflect.NewAt(field.Type, unsafe.Pointer(fieldp))
				f.format(s, c, reflect.Indirect(fieldptr))
			}
		}
		f.depth--
		if f.verbose && f.pretty {
			writeComma(s)
			writeNewline(s)
			writeFullIndent(s, f.depth)
		}
		writeRightcurly(s)
	default:
		fmt.Fprintf(s, reconstructFlags(s, 'v'), val.Interface())
	}
}

func (f *formatter) sep(s fmt.State) {
	if f.verbose {
		writeComma(s)
		if f.pretty {
			writeNewline(s)
			writeFullIndent(s, f.depth)
			return
		}
	}
	writeSpace(s)
}

func reconstructFlags(s fmt.State, c rune) string {
	flags := make([]rune, 0, 7)
	flags = append(flags, '%')
	flags = addFlagRune(flags, s, '+')
	flags = addFlagRune(flags, s, '-')
	flags = addFlagRune(flags, s, '#')
	flags = addFlagRune(flags, s, ' ')
	flags = addFlagRune(flags, s, '0')
	flags = append(flags, c)
	return string(flags)
}

func addFlagRune(q []rune, s fmt.State, r rune) []rune {
	if s.Flag(int(r)) {
		return append(q, r)
	}
	return q
}

func Printf(format string, v ...interface{}) {
	Fprintf(os.Stdout, format, v...)
}
func Sprintf(format string, v ...interface{}) string {
	buf := new(bytes.Buffer)
	Fprintf(buf, format, v...)
	return buf.String()
}
func Fprintf(w io.Writer, format string, v ...interface{}) {
	_v := make([]interface{}, len(v))
	for i := range v {
		_v[i] = NewFormatter(v[i])
	}
	fmt.Fprintf(w, format, _v...)
}
