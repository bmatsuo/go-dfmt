package dfmt

import (
	"fmt"
	"testing"
)

type testStruct struct {
	X *testStruct
	y interface{}
}

// %v and %+v as well as %#v and %#V output the same strings when desired
func TestCompatability(t *testing.T) {
	for i, test := range []interface{}{
		nil,
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
		std := fmt.Sprintf("%v", test)
		deep := fmt.Sprintf("%v", Formatter(debugging, test))
		if std != deep {
			t.Errorf("[%d] # mismatch %q != %q ", i, std, deep)
		}
		std = fmt.Sprintf("%#v", test)
		deep = fmt.Sprintf("%#v", Formatter(debugging, test))
		if std != deep {
			t.Errorf("[%d] # mismatch %q != %q ", i, std, deep)
		}
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
