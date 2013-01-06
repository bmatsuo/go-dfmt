package dfmt

import (
	"fmt"
	"errors"
	"testing"
)

type testStruct struct {
	X *testStruct
	y interface{}
}

// %v and %+v as well as %#v and %#V output the same strings when desired
func TestCompatability(t *testing.T) {
	testflags := func(flags string, i int, test interface{}) {
		std := fmt.Sprintf("%"+flags+"v", test)
		deep := fmt.Sprintf("%"+flags+"v", Formatter(debugging, test))
		if std != deep {
			t.Errorf("[%d] "+flags+" mismatch %q != %q ", i, std, deep)
		}
	}
	for i, test := range []interface{}{
		nil,
		errors.New("foobar"),
		(*testStruct)(nil),
		"abc",
		12 + 4i,
		[]string{"a", "b"},
		map[string]int64{"a": 1}, // map iteration order consistant
		map[string]interface{}{"a": &testStruct{}},
		[]interface{}{&testStruct{}},
		&struct{ A, b int }{1, 2},
		struct{ A, b int }{1, 2},
		testStruct{nil, map[string]int{"abc": 2}},
		&testStruct{nil, map[string]int{"abc": 2}},
		testStruct{&testStruct{}, map[string]int{"abc": 2}},
		&testStruct{&testStruct{}, map[string]int{"abc": 2}},
	} {
		testflags("", i, test)

		testflags("#", i, test)

		testflags("+", i, test)
		testflags("#+", i, test)

		testflags("-", i, test)
		testflags("#-", i, test)
		testflags("+-", i, test)
		testflags("#+-", i, test)

		testflags("0", i, test)
		testflags("-0", i, test)
		testflags("#-0", i, test)
		testflags("+-0", i, test)
		testflags("#+-0", i, test)

		testflags(" ", i, test)
		testflags("0 ", i, test)
		testflags("-0 ", i, test)
		testflags("#-0 ", i, test)
		testflags("+-0 ", i, test)
		testflags("#+-0 ", i, test)
	}
}

func TestDeepfmt(t *testing.T) {
	/*
		printf(Deep, "%v\n", nil)
		printf(Deep, "%#v\n", (*testStruct)(nil))
		printf(Deep, "%v\n", "abc")
		printf(Deep, "%#v\n", "abc")
		printf(Deep, "%v\n", 12+4i)
		printf(Deep, "%#v\n", 12+4i)
		printf(Deep, "%v\n", []string{"a", "b"})
		printf(Deep, "%#v\n", []string{"a", "b"})
		printf(Deep|Pretty, "%#v\n", []string{"a", "b"})
		printf(Deep, "%v\n", map[string]int64{"a": 1, "b": 2})
		printf(Deep, "%#v\n", map[string]int64{"a": 1, "b": 2})
		printf(Deep|Pretty, "%#v\n", map[string]int64{"a": 1, "b": 2})
		printf(Deep, "%v\n", &struct{ A, b int }{1, 2})
		printf(Deep, "%#v\n", &struct{ A, b int }{1, 2})
		printf(Deep|Pretty, "%#v\n", &struct{ A, b int }{1, 2})
		printf(Deep, "%v\n", struct{ A, b int }{1, 2})
		printf(Deep, "%#v\n", struct{ A, b int }{1, 2})
		printf(Deep|Pretty, "%#v\n", struct{ A, b int }{1, 2})
		printf(Deep, "%#v\n", testStruct{nil, map[string]int{"abc": 2}})
		printf(Deep, "%#v\n", &testStruct{nil, map[string]int{"abc": 2}})
		printf(Deep|debugging, "%#v\n", &testStruct{&testStruct{}, map[string]int{"abc": 2}})

		var x ****testStruct
		x = new(***testStruct)
		*x = new(**testStruct)
		**x = new(*testStruct)
		printf(Deep, "%v\n", x)
		printf(Deep, "%#v\n", x)
	*/
}
