package dfmt

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

type FormatMode uint

const (
	noop     FormatMode = 1 << iota
	Deep                // follow pointers
	Pretty              // pretty print
	NoIface             // ignore fmt interfaces
	NoComply            // use noncompliant flags TODO
	debug
)

type formatter struct {
	depth   int
	verbose bool
	extra   bool
	space   bool
	deep    bool
	pretty  bool
	ifaceok bool
	debug   bool
	v       interface{}
}

// Wrap v in a fmt.Formatter that overrides %v and applies FormatModes.
func Formatter(mode FormatMode, v interface{}) fmt.Formatter {
	return &formatter{
		deep:    mode&Deep != 0,
		pretty:  mode&Pretty != 0,
		ifaceok: mode&NoIface == 0,
		debug:   mode&debug != 0,
		v:       v,
	}
}

func (f *formatter) Format(s fmt.State, c rune) {
	if c != 'v' || (!f.deep && !f.pretty && f.ifaceok && !f.debug) {
		// default printf
		newfmt := reconstructFlags(s, c)
		fmt.Fprintf(s, newfmt, f.v)
		return
	}
	if f.v == nil {
		writeNilangle(s)
		return
	}

	f.verbose = s.Flag('#')
	f.extra = s.Flag('+')
	f.space = s.Flag(' ')
	f.format(s, c, reflect.ValueOf(f.v))
}

func (f *formatter) format(s fmt.State, c rune, val reflect.Value) {
	if f.ifaceok {
		// check fmt interface overrides
		v := val.Interface()
		if formatter, ok := v.(fmt.Formatter); ok {
			formatter.Format(s, c)
			return
		}
		if !f.verbose {
			switch v.(type) {
			case fmt.Stringer:
				fmt.Fprint(s, v.(fmt.Stringer).String())
				return
			case error:
				fmt.Fprint(s, v.(error).Error())
				return
			}
		} else if gs, ok := v.(fmt.GoStringer); ok {
			fmt.Fprint(s, gs.GoString())
			return
		}
	}

	switch val.Kind() {
	case reflect.Interface:
		f.formatInterface(s, c, val)
	case reflect.Ptr:
		f.formatPtr(s, c, val)
	case reflect.Array, reflect.Slice:
		f.formatArray(s, c, val)
	case reflect.Map:
		f.formatMap(s, c, val)
	case reflect.Struct:
		f.formatStruct(s, c, val)
	default:
		fmt.Fprintf(s, reconstructFlags(s, 'v'), val.Interface())
	}
}
func (f *formatter) formatInterface(s fmt.State, c rune, val reflect.Value) {
	if val.IsNil() {
		if f.verbose {
			writeType(s, val)
			writeLeftparen(s)
			writeNil(s)
			writeRightparen(s)
		} else {
			writeNilangle(s)
		}
		return
	}
	f.format(s, c, val.Elem())
}
func (f *formatter) formatPtr(s fmt.State, c rune, val reflect.Value) {
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
	if !f.deep && f.depth > 0 {
		if f.verbose {
			writeLeftparen(s)
			writeType(s, val)
			writeRightparen(s)
			writeLeftparen(s)
			if f.space {
				writeSpace(s)
			}
			fmt.Fprintf(s, "%p", val.Interface())
			writeRightparen(s)
		} else {
			if f.space {
				writeSpace(s)
			}
			fmt.Fprintf(s, "%p", val.Interface())
		}
		return
	}
	writePtr(s)
	f.format(s, c, val.Elem())
}
func (f *formatter) formatArray(s fmt.State, c rune, val reflect.Value) {
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
}
func (f *formatter) formatMap(s fmt.State, c rune, val reflect.Value) {
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
}
func (f *formatter) formatStruct(s fmt.State, c rune, val reflect.Value) {
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
		if f.verbose || f.extra {
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

// for test/debug
func printf(mode FormatMode, format string, v ...interface{}) {
	fprintf(os.Stdout, mode, format, v...)
}
func sprintf(mode FormatMode, format string, v ...interface{}) string {
	buf := new(bytes.Buffer)
	fprintf(buf, mode, format, v...)
	return buf.String()
}
func fprintf(w io.Writer, mode FormatMode, format string, v ...interface{}) {
	_v := make([]interface{}, len(v))
	for i := range v {
		_v[i] = Formatter(mode, v[i])
	}
	fmt.Fprintf(w, format, _v...)
}
