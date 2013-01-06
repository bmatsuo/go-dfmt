/*
A drop-in replacement for the standard library fmt package using dfmt.Formatter.

This package is not performant. Use it during development but keep fmt in
production code.
*/
package fmt

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/bmatsuo/go-dfmt"
)

func Errorf(format string, a ...interface{}) error {
	for i := range a {
		a[i] = dfmt.Formatter(a[i])
	}
	return fmt.Errorf(format, a)
}
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	for i := range a {
		a[i] = dfmt.Formatter(a[i])
	}
	return fmt.Fprintf(w, format, a)
}
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, format, a...)
}
func Sprintf(format string, a ...interface{}) string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, format, a...)
	return b.String()
}

func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(os.Stdout, a...)
}
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, a...)
}

func Fscan(r io.Reader, a ...interface{}) (n int, err error) {
	return Fscan(r, a...)
}
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error) {
	return Fscanf(r, format, a...)
}
func Fscanln(r io.Reader, a ...interface{}) (n int, err error) {
	return Fscanln(r, a...)
}

func Print(a ...interface{}) (n int, err error) {
	return fmt.Print(a...)
}
func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func Scan(a ...interface{}) (n int, err error) {
	return fmt.Scan(a...)
}
func Scanf(format string, a ...interface{}) (n int, err error) {
	return fmt.Scanf(format, a...)
}
func Scanln(a ...interface{}) (n int, err error) {
	return fmt.Scanln(a...)
}

func Sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}
func Sprintln(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func Sscan(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscan(str, a...)
}
func Sscanf(str string, format string, a ...interface{}) (n int, err error) {
	return fmt.Sscanf(str, format, a...)
}
func Sscanln(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscanln(str, a...)
}

type Formatter fmt.Formatter
type GoStringer fmt.GoStringer
type ScanState fmt.ScanState
type Scanner fmt.Scanner
type State fmt.State
type Stringer fmt.Stringer
